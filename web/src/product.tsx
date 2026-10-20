import React,{useEffect,useRef,useState} from 'react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import {Maps} from './generated/Maps';

type Place={id?:string,key?:string,name:string,address:string,lat:number,lon:number,list?:string,notes?:string};
const point=(f:any,i:number):Place=>({key:String(i),name:f.properties.name||f.properties.street||'Place',address:[f.properties.housenumber,f.properties.street,f.properties.city,f.properties.state,f.properties.country].filter(Boolean).join(', '),lon:f.geometry.coordinates[0],lat:f.geometry.coordinates[1]});
const fmtDistance=(m:number)=>m<1000?Math.round(m)+' m':(m/1000).toFixed(1)+' km';
function MapView({mapRef,onPoint,tiles,onError,onReady}:any){
 const node=useRef<HTMLDivElement>(null),click=useRef(onPoint);click.current=onPoint;
 useEffect(()=>{const map=L.map(node.current!,{zoomControl:false}).setView([51.505,-0.09],12);mapRef.current=map;
 L.control.zoom({position:'bottomright'}).addTo(map);
 const layer=L.tileLayer(tiles,{maxZoom:19,attribution:'© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',updateWhenIdle:true}).addTo(map);
 let reported=false;layer.on('tileerror',()=>{if(!reported){onError('Map tiles could not load. Check your connection; your saved places are still available.');reported=true}});
 map.on('click',(e:any)=>click.current(e.latlng));const obs=new ResizeObserver(()=>map.invalidateSize());obs.observe(node.current!);onReady();
 return()=>{obs.disconnect();map.remove();mapRef.current=null}},[tiles]);
 return React.createElement('div',{ref:node,className:'map-canvas',role:'region','aria-label':'Interactive map'});
}
export function Product(p:any){
 const [tab,setTab]=useState('explore'),[config,setConfig]=useState<any>(null),[results,setResults]=useState<Place[]>([]),[places,setPlaces]=useState<Place[]>([]),[selected,setSelected]=useState<Place|null>(null),[stops,setStops]=useState<Place[]>([]),[routes,setRoutes]=useState<any[]>([]),[routeIndex,setRouteIndex]=useState(0),[measuring,setMeasuring]=useState(false),[measurePoints,setMeasurePoints]=useState<any[]>([]),[navigating,setNavigating]=useState(false);
 const [mapGeneration,setMapGeneration]=useState(0);
 const mapRef=useRef<L.Map|null>(null),pins=useRef<L.LayerGroup|null>(null),line=useRef<L.Polyline|null>(null),measurementLine=useRef<L.Polyline|null>(null),watch=useRef<number|null>(null),locationMarker=useRef<L.CircleMarker|null>(null),alive=useRef(true),selection=useRef(0),f=p.fields;
 const loadPlaces=async()=>{const x=await p.request('/api/places');if(alive.current)setPlaces(x)};
 useEffect(()=>{alive.current=true;p.request('/api/config').then((v:any)=>alive.current&&setConfig(v)).catch((e:any)=>p.setError(e.message));loadPlaces().catch((e:any)=>p.setError(e.message));return()=>{alive.current=false;if(watch.current!==null)navigator.geolocation.clearWatch(watch.current);speechSynthesis?.cancel()}},[]);
 useEffect(()=>{if(!mapRef.current)return;const map=mapRef.current;pins.current?.remove();pins.current=L.layerGroup().addTo(map);
 [...(tab==='saved'?places:results),...stops,...(selected?[selected]:[])].forEach(x=>{const m=L.circleMarker([x.lat,x.lon],{radius:8,color:'#fff',weight:3,fillColor:'#4656df',fillOpacity:1}).addTo(pins.current!);m.bindTooltip(document.createTextNode(x.name) as any);m.on('click',()=>choose(x))});
 },[places,results,stops,selected,tab,config,mapGeneration]);
 useEffect(()=>{const map=mapRef.current;if(!map)return;line.current?.remove();const r=routes[routeIndex];if(r){line.current=L.polyline(r.geometry.coordinates.map((c:number[])=>[c[1],c[0]]),{color:'#485ce0',weight:6}).addTo(map);map.fitBounds(line.current.getBounds(),{padding:[40,40]})}else if(selected){map.setView([selected.lat,selected.lon],14)}},[routes,routeIndex,mapGeneration]);
 useEffect(()=>{const map=mapRef.current;if(!map)return;measurementLine.current?.remove();if(measurePoints.length)measurementLine.current=L.polyline(measurePoints,{color:'#13aa93',dashArray:'6 7'}).addTo(map)},[measurePoints,mapGeneration]);
 useEffect(()=>{if(!config)return;const q=new URLSearchParams(location.hash.slice(1));const lat=Number(q.get('lat')),lon=Number(q.get('lon'));if(q.has('lat')&&q.has('lon')&&Number.isFinite(lat)&&Math.abs(lat)<=90&&Number.isFinite(lon)&&Math.abs(lon)<=180){const x={name:q.get('name')||'Shared place',address:'',lat,lon};choose(x)}},[config]);
 const choose=(x:Place)=>{selection.current++;setSelected(x);p.setField('list',x.list||'Favourites');p.setField('notes',x.notes||'');mapRef.current?.setView([x.lat,x.lon],Math.max(mapRef.current.getZoom(),14))};
 const pick=async(x:any)=>{if(measuring){setMeasurePoints(v=>[...v,[x.lat,x.lng]]);return}const pos={lat:x.lat,lon:x.lng,name:'Dropped pin',address:''};choose(pos);const request=selection.current;try{const r=await p.request('/api/reverse?lat='+x.lat+'&lon='+x.lng);if(alive.current&&selection.current===request&&r.features?.length)choose({...point(r.features[0],0),lat:x.lat,lon:x.lng})}catch(e:any){if(selection.current===request)p.setError(e.message)}};
 const locate=(continuous=false)=>{if(!navigator.geolocation)throw Error('Location is unavailable in this browser');
 const success=(pos:GeolocationPosition)=>{if(!alive.current||!mapRef.current)return;const ll:L.LatLngExpression=[pos.coords.latitude,pos.coords.longitude];locationMarker.current?.remove();locationMarker.current=L.circleMarker(ll,{color:'#1479ee',radius:7,fillOpacity:1}).addTo(mapRef.current);mapRef.current.setView(ll,16);if(!continuous)choose({lat:pos.coords.latitude,lon:pos.coords.longitude,name:'My location',address:''})};
 const fail=(e:GeolocationPositionError)=>{setNavigating(false);p.setError(e.code===1?'Location permission was declined. You can still choose a starting point on the map.':'Could not locate this device. Choose your start on the map.')};
 if(continuous){watch.current=navigator.geolocation.watchPosition(success,fail,{enableHighAccuracy:true,maximumAge:5000});setNavigating(true)}else navigator.geolocation.getCurrentPosition(success,fail,{enableHighAccuracy:true,timeout:15000})};
 const action=async(a:string)=>{if(['apps','closeApps','logout','theme'].includes(a))return p.common(a);p.setError('');const id=p.selected.current;try{switch(a){
 case'explore':setTab('explore');break;case'saved':setTab('saved');break;case'directions':setTab('directions');break;
 case'search':p.setBusy('Finding places…');{const r=await p.request('/api/search?q='+encodeURIComponent(f.query||''));const v=(r.features||[]).map(point);setResults(v);if(v.length){choose(v[0]);setTab('explore')}else p.setNotice('No places found. Try a town, address or landmark.')}break;
 case'place':{const x=results.find(x=>x.key===id);if(x)choose(x);break}case'savedPlace':{const x=places.find(x=>x.id===id);if(x)choose(x);break}
 case'deselect':selection.current++;setSelected(null);break;case'start':if(selected){setStops(s=>[selected,...s.slice(1)]);setRoutes([]);setTab('directions')}break;case'stop':if(selected){if(stops.length>=8)throw Error('Use up to eight stops per route');setStops(s=>[...s,selected]);setRoutes([]);setTab('directions')}break;
 case'removeStop':setStops(s=>s.filter((_,i)=>String(i)!==id));setRoutes([]);break;case'clearStops':setStops([]);setRoutes([]);break;
 case'route':p.setBusy('Calculating route…');{const r=await p.request('/api/route',{method:'POST',body:JSON.stringify({mode:f.mode||'car',points:stops.map(x=>[x.lon,x.lat])})});if(r.code!=='Ok'||!r.routes?.length)throw Error('No route connects these stops. Try a different travel mode or move a stop onto a road.');setRoutes(r.routes);setRouteIndex(0)}break;
 case'chooseRoute':setRouteIndex(Number(id));break;
 case'save':if(selected){const x=await p.request('/api/places',{method:'POST',body:JSON.stringify({...selected,id:selected.id||'',key:undefined,list:f.list||'Favourites',notes:f.notes||''})});setSelected(x);await loadPlaces();p.setNotice('Place saved')}break;
 case'deletePlace':await p.request('/api/places/'+id,{method:'DELETE'});await loadPlaces();if(selected?.id===id)setSelected(null);break;
 case'share':if(selected){const u=new URL(location.href);u.hash=new URLSearchParams({lat:String(selected.lat),lon:String(selected.lon),name:selected.name}).toString();await navigator.clipboard.writeText(u.href);p.setNotice('Place link copied. Only its name and coordinates are shared.')}break;
 case'locate':locate();break;case'measure':setMeasuring(!measuring);setMeasurePoints([]);break;
 case'navigate':if(navigating){if(watch.current!==null)navigator.geolocation.clearWatch(watch.current);watch.current=null;setNavigating(false)}else locate(true);break;
 case'exportPlaces':await p.download(new Blob([JSON.stringify({type:'FeatureCollection',features:places.map(x=>({type:'Feature',properties:{name:x.name,address:x.address,list:x.list,notes:x.notes},geometry:{type:'Point',coordinates:[x.lon,x.lat]}}))},null,2)],{type:'application/geo+json'}),'my-places.geojson');break;
 case'exportRoute':{const r=routes[routeIndex];if(r)await p.download(new Blob(['<?xml version="1.0"?><gpx version="1.1" creator="My Maps" xmlns="http://www.topografix.com/GPX/1/1"><trk><name>My Maps route</name><trkseg>'+r.geometry.coordinates.map((c:number[])=>`<trkpt lat="${c[1]}" lon="${c[0]}"/>`).join('')+'</trkseg></trk></gpx>'],{type:'application/gpx+xml'}),'route.gpx');break}
 }}catch(e:any){p.setError(e.message)}finally{p.setBusy('')}};
 if(p.route==='#settings')return p.render(p.settings(),{},action);
 const r=routes[routeIndex];const steps=r?.legs.flatMap((l:any)=>l.steps).map((s:any)=>({instruction:(s.maneuver.type==='depart'?'Start':s.maneuver.type==='arrive'?'Arrive':s.maneuver.type.replaceAll('_',' '))+(s.maneuver.modifier?' '+s.maneuver.modifier:'')+(s.name?' onto '+s.name:''),distance:fmtDistance(s.distance)}))||[];
 const data={exploring:tab==='explore',saving:tab==='saved',routing:tab==='directions',exploreVariant:tab==='explore'?'active':'secondary',savedVariant:tab==='saved'?'active':'secondary',routeVariant:tab==='directions'?'active':'secondary',results,places,noPlaces:!places.length,selected:selected?{...selected,coordinates:selected.lat.toFixed(5)+', '+selected.lon.toFixed(5)}:null,stops:stops.map((s,i)=>({key:String(i),label:(i+1)+'. '+s.name})),notEnoughStops:stops.length<2,modes:[{value:'car',label:'Driving'},{value:'bike',label:'Cycling'},{value:'foot',label:'Walking'}],hasRoutes:routes.length>0,routes:routes.map((r,i)=>({key:String(i),variant:i===routeIndex?'active':'secondary',summary:Math.round(r.duration/60)+' min · '+fmtDistance(r.distance)})),steps,navigationLabel:navigating?'Stop following':'Follow my location',measureLabel:measuring?'Finish measuring':'Measure distance',measurement:measuring?'Click the map: '+fmtDistance(measurePoints.reduce((n,x,i)=>i?n+L.latLng(measurePoints[i-1]).distanceTo(x):0,0)):''};
 return p.render(React.createElement(Maps,{data,onAction:action}),{},action,{map:config?React.createElement(MapView,{mapRef,onPoint:pick,tiles:config.tiles,onError:p.setError,onReady:()=>setMapGeneration(v=>v+1)}):null});
}
