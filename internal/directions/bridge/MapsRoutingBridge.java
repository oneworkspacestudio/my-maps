import java.io.*;
import java.net.*;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicReference;
import com.sun.net.httpserver.*;
import btools.mapaccess.PhysicalFile;
import btools.router.*;
import btools.server.ServiceContext;
import btools.server.request.ServerHandler;

/** Private Maps adapter. Deliberately exposes no BRouter profile-upload,
 * suspects/history, file, or unauthenticated routing endpoints. */
public final class MapsRoutingBridge {
  private static final Set<String> PROFILES=Set.of("car-fast","car-eco","trekking","fastbike","gravel","mtb","shortest","hiking-mountain");
  private static final Set<String> OPTIONS=Set.of("profile:avoid_toll","profile:avoid_motorways","profile:avoid_unpaved","profile:recup_efficiency");
  private static final AtomicReference<RoutingEngine> active=new AtomicReference<>();
  private static final AtomicReference<String> activeId=new AtomicReference<>();
  private static final AtomicReference<String> cancelledId=new AtomicReference<>();
  private static void reply(HttpExchange x,int status,String body) throws IOException {
    byte[] bytes=body.getBytes(StandardCharsets.UTF_8);
    x.getResponseHeaders().set("Content-Type","application/json; charset=utf-8");
    x.getResponseHeaders().set("Cache-Control","no-store");
    x.sendResponseHeaders(status,bytes.length);
    try(OutputStream out=x.getResponseBody()){out.write(bytes);}
  }
  public static void main(String[] args) throws Exception {
    if(args.length!=2)throw new IllegalArgumentException("Expected private data directories");
    String token=new BufferedReader(new InputStreamReader(System.in,StandardCharsets.UTF_8)).readLine();
    if(token==null||!token.matches("[a-f0-9]{64}"))throw new IllegalArgumentException("Missing process token");
    PrintStream control=System.out;
    // The upstream server logs precise route URLs and has a recent-waypoint
    // endpoint. Neither is used. Suppress all underlying engine diagnostics.
    System.setOut(new PrintStream(OutputStream.nullOutputStream()));
    System.setErr(new PrintStream(OutputStream.nullOutputStream()));
    File segments=new File(args[0]);
    File[] files=segments.listFiles((dir,name)->name.endsWith(".rd5"));
    if(files==null||files.length==0)throw new IOException("Missing segments");
    for(File file:files){String error=PhysicalFile.checkFileIntegrity(file);if(error!=null)throw new IOException("Segment integrity failed");}
    ServiceContext service=new ServiceContext();service.segmentDir=segments;service.profileDir=args[1];
    service.customProfileDir=args[1];service.sharedProfileDir=args[1];
    System.setProperty("profileBaseDir",service.profileDir);
    HttpServer server=HttpServer.create(new InetSocketAddress(InetAddress.getByName("127.0.0.1"),0),4);
    server.setExecutor(new ThreadPoolExecutor(2,2,0,TimeUnit.SECONDS,new ArrayBlockingQueue<>(4),new ThreadPoolExecutor.AbortPolicy()));
    server.createContext("/",x->{
      try{
        if(!MessageDigest.isEqual(("Bearer "+token).getBytes(StandardCharsets.UTF_8),Objects.toString(x.getRequestHeaders().getFirst("Authorization"),"").getBytes(StandardCharsets.UTF_8))){reply(x,403,"{}");return;}
        String path=x.getRequestURI().getPath();
        if(path.equals("/healthz")&&x.getRequestMethod().equals("GET")){reply(x,200,"{\"status\":\"ok\"}");return;}
        if(!x.getRequestMethod().equals("POST")||(!path.equals("/route")&&!path.equals("/cancel"))){reply(x,404,"{}");return;}
        byte[] body=x.getRequestBody().readNBytes(16385);if(body.length>16384){reply(x,413,"{}");return;}
        Map<String,String> params=new HashMap<>();
        for(String part:new String(body,StandardCharsets.UTF_8).split("&")){
          String[] pair=part.split("=",2);if(pair.length!=2)throw new IllegalArgumentException();
          String key=URLDecoder.decode(pair[0],StandardCharsets.UTF_8),value=URLDecoder.decode(pair[1],StandardCharsets.UTF_8);
          if(params.put(key,value)!=null)throw new IllegalArgumentException();
        }
        String id=params.remove("requestId");if(id==null||!id.matches("[a-f0-9]{32}"))throw new IllegalArgumentException();
        if(path.equals("/cancel")){RoutingEngine e=active.get();if(id.equals(activeId.get())){cancelledId.set(id);if(e!=null)e.terminate();}reply(x,200,"{}");return;}
        if(!PROFILES.contains(params.get("profile")))throw new IllegalArgumentException();
        for(String key:params.keySet())if(!Set.of("profile","lonlats","alternativeidx").contains(key)&&!OPTIONS.contains(key))throw new IllegalArgumentException();
        for(String key:OPTIONS)if(params.containsKey(key)&&!params.get(key).matches("0|1|0\\.[0-9]{1,2}"))throw new IllegalArgumentException();
        String[] points=Objects.requireNonNull(params.get("lonlats")).split("\\|");if(points.length<2||points.length>8)throw new IllegalArgumentException();
        for(String point:points){String[] xy=point.split(",");if(xy.length!=2)throw new IllegalArgumentException();double lon=Double.parseDouble(xy[0]),lat=Double.parseDouble(xy[1]);if(!Double.isFinite(lon)||!Double.isFinite(lat)||Math.abs(lon)>180||Math.abs(lat)>85)throw new IllegalArgumentException();}
        if(!Objects.toString(params.get("alternativeidx"),"0").matches("[012]"))throw new IllegalArgumentException();
        if(!activeId.compareAndSet(null,id)){reply(x,429,"{\"error\":\"busy\"}");return;}
        try{
          params.put("format","geojson");params.put("timode","2");
          ServerHandler handler=new ServerHandler(service,params);RoutingContext rc=handler.readRoutingContext();rc.memoryclass=256;
          RoutingParamCollector collector=new RoutingParamCollector();List<OsmNodeNamed> waypoints=collector.getWayPointList(params.get("lonlats"));
          collector.setParams(rc,waypoints,params);
          RoutingEngine engine=new RoutingEngine(null,null,service.segmentDir,waypoints,rc,0);engine.quite=true;active.set(engine);if(id.equals(cancelledId.get()))engine.terminate();engine.doRun(90_000);
          if(engine.getErrorMessage()!=null||engine.getFoundTrack()==null){reply(x,422,"{\"error\":\"no_route\"}");return;}
          String result=handler.formatTrack(engine.getFoundTrack());if(result.length()>8_000_000){reply(x,413,"{}");return;}
          reply(x,200,result);
        }finally{active.set(null);activeId.set(null);cancelledId.set(null);}
      }catch(Throwable e){try{reply(x,400,"{\"error\":\"route_failed\"}");}catch(IOException ignored){}}
      finally{x.close();}
    });
    server.start();control.println("READY "+server.getAddress().getPort());control.flush();
  }
}
