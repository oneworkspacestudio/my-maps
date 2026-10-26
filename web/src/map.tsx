import React,{useEffect,useRef,useState} from 'react';
import {Map as GLMap, Marker, NavigationControl, AttributionControl, setWorkerUrl, type MapMouseEvent} from 'maplibre-gl';
import workerURL from 'maplibre-gl/dist/maplibre-gl-worker.mjs?worker&url';
import 'maplibre-gl/dist/maplibre-gl.css';
import {RasterMap} from './raster-map';
import {vectorController} from './map-controller';

setWorkerUrl(workerURL);
export type Place={id?:string,key?:string,name:string,address:string,lat:number,lon:number,list?:string,notes?:string,category?:string,hours?:string,phone?:string,website?:string,osm?:string};
export const distance=(a:number[],b:number[])=>{const rad=Math.PI/180,lat=(b[1]-a[1])*rad,lon=(b[0]-a[0])*rad;return 12742000*Math.asin(Math.sqrt(Math.sin(lat/2)**2+Math.cos(a[1]*rad)*Math.cos(b[1]*rad)*Math.sin(lon/2)**2))};
export const fmtDistance=(m:number)=>m<1000?Math.round(m)+' m':(m/1000).toFixed(1)+' km';
export function MapView(p:any){
 const [failed,setFailed]=useState(false),phone=useRef(matchMedia('(any-pointer: coarse)').matches||/Android|iPhone|iPad|iPod/i.test(navigator.userAgent));
 const raster=failed||(phone.current&&p.mapType!=='terrain'&&!p.buildings);
 const fallback=()=>{setFailed(true);p.onFallback()};
 return React.createElement(raster?RasterMap:VectorMap,{...p,onFallback:fallback});
}
function VectorMap(p:any){
 const defaultStyle=p.config.customTiles?{version:8,sources:{custom:{type:'raster',tiles:[p.config.tiles],tileSize:256,attribution:p.config.attribution}},layers:[{id:'custom',type:'raster',source:'custom'}]}:'https://tiles.openfreemap.org/styles/liberty';
 const node=useRef<HTMLDivElement>(null),native=useRef<GLMap|null>(null),latest=useRef(p),previousType=useRef('default'),vectorZoom=useRef(12);latest.current=p;
 useEffect(()=>{
  let reported=false;
  let map:GLMap;
  try{map=new GLMap({container:node.current!,style:defaultStyle as any,center:[-.09,51.505],zoom:12,...(p.selected?{center:[p.selected.lon,p.selected.lat] as [number,number],zoom:15}:{}),...p.camera.current,attributionControl:false,transformRequest:(url:string)=>url.startsWith('https://')?{url:p.base+'/api/map-resource?url='+encodeURIComponent(url)}:{url}})}catch{latest.current.onFallback();return}
  native.current=map;const controller=vectorController(map);p.mapRef.current=controller;
  // A style may load while its worker/tiles never render. Recover that blank
  // canvas as well as GPU creation/context failures; never wait indefinitely.
  let deadline=Date.now()+12000;
  const watchdog=setInterval(()=>{if(document.hidden){deadline=Date.now()+12000;return}const style=map.getStyle();if(map.isStyleLoaded()&&(map.queryRenderedFeatures().length||(Object.values(style.sources).some(s=>s.type==='raster')&&map.areTilesLoaded()))){deadline=Date.now()+12000;return}if(Date.now()>deadline){clearInterval(watchdog);latest.current.onFallback()}},1000);
  const lost=()=>latest.current.onFallback();map.on('webglcontextlost',lost);
  map.addControl(new NavigationControl({showCompass:false}),'bottom-right');
  map.addControl(new AttributionControl({compact:true}),'bottom-left');
  map.on('error',()=>{if(!reported){reported=true;latest.current.onError('Some map data could not load. Check your connection or try again.')}});
  map.on('click',(e:MapMouseEvent)=>{
   const feature=map.queryRenderedFeatures(e.point).find((f:any)=>f.layer.type==='symbol'&&/poi|place|transport/.test(f.layer.id)&&f.properties?.name);
   latest.current.onPoint({lat:e.lngLat.lat,lon:e.lngLat.lng,name:feature?.properties?.name});
  });
  map.on('moveend',()=>latest.current.onMove());
  map.on('style.load',()=>{
   reported=false;
   if(latest.current.mapType!=='satellite'){
    for(const layer of map.getStyle().layers){
     if(layer.type==='fill'&&/water$/.test(layer.id))map.setPaintProperty(layer.id,'fill-color','#95d9e8');
     if(layer.type==='fill'&&/park|landcover_wood|landuse_grass/.test(layer.id))map.setPaintProperty(layer.id,'fill-color','#bfe9cc');
     if(layer.type==='line'&&(layer as any)['source-layer']==='transportation'&&!/rail|ferry|aerialway|waterway/.test(layer.id))map.setPaintProperty(layer.id,'line-color',/casing|outline/.test(layer.id)?'#c3cbd1':/motorway|trunk/.test(layer.id)?'#f3eac8':'#ffffff');
    }
    const source=Object.entries(map.getStyle().sources).find(([,s])=>s.type==='vector')?.[0];
    const before=map.getStyle().layers.find(l=>l.type==='symbol')?.id;
    if(source&&!map.getLayer('raised-buildings'))map.addLayer({id:'raised-buildings',type:'fill-extrusion',source,'source-layer':'building',minzoom:14,layout:{visibility:latest.current.buildings?'visible':'none'},paint:{'fill-extrusion-color':'#c5cad0','fill-extrusion-height':['coalesce',['get','render_height'],6],'fill-extrusion-base':['coalesce',['get','render_min_height'],0],'fill-extrusion-opacity':.88}},before);
   }
   if(latest.current.mapType==='terrain'){
    map.addSource('terrain-dem',{type:'raster-dem',url:'https://tiles.mapterhorn.com/tilejson.json'});
    map.setTerrain({source:'terrain-dem',exaggeration:1.2});
    map.addLayer({id:'terrain-shading',type:'hillshade',source:'terrain-dem',paint:{'hillshade-exaggeration':.25}},map.getStyle().layers.find(l=>l.type==='symbol')?.id);
   }
   latest.current.onReady();
  });
  const resize=new ResizeObserver(()=>map.resize());resize.observe(node.current!);
  return()=>{clearInterval(watchdog);resize.disconnect();latest.current.onMove();map.off('webglcontextlost',lost);map.remove();native.current=null;if(p.mapRef.current===controller)p.mapRef.current=null};
 },[]);
 useEffect(()=>{
  const map=native.current;if(!map)return;
  if(p.mapType==='satellite'){vectorZoom.current=map.getZoom();map.setStyle({version:8,sources:{imagery:{type:'raster',tiles:['https://gibs.earthdata.nasa.gov/wmts/epsg3857/best/BlueMarble_ShadedRelief/default/GoogleMapsCompatible_Level8/{z}/{y}/{x}.jpeg'],tileSize:256,maxzoom:8,attribution:'Imagery: NASA Blue Marble (500 m)'}},layers:[{id:'imagery',type:'raster',source:'imagery'}]});map.easeTo({zoom:Math.min(8,map.getZoom()),pitch:0})}
  else map.setStyle((p.mapType==='default'?defaultStyle:'https://tiles.openfreemap.org/styles/liberty') as any);
  if(p.mapType!=='satellite')map.easeTo({pitch:p.mapType==='terrain'?45:0,...(previousType.current==='satellite'?{zoom:vectorZoom.current}:{})});previousType.current=p.mapType;
 },[p.mapType]);
 useEffect(()=>{const map=native.current;if(map?.getLayer('raised-buildings')){map.setLayoutProperty('raised-buildings','visibility',p.buildings?'visible':'none');map.easeTo({pitch:p.buildings?50:p.mapType==='terrain'?45:0,zoom:p.buildings?Math.max(15,map.getZoom()):map.getZoom()})}},[p.buildings,p.generation]);
 useEffect(()=>{
  const map=native.current;if(!map)return;const markers:Marker[]=[];
  const unique=new Map<string,Place>();for(const x of [...(p.selected?[p.selected]:[]),...p.stops,...p.places])if(!unique.has(x.lat+','+x.lon))unique.set(x.lat+','+x.lon,x);
  const occupied:{x:number,y:number}[]=[];for(const x of unique.values()){const point=map.project([x.lon,x.lat]);if(x!==p.selected&&occupied.some(v=>Math.hypot(v.x-point.x,v.y-point.y)<32))continue;occupied.push(point);
   const el=document.createElement('button');el.className='place-pin'+(x===p.selected?' place-pin--selected':'');el.type='button';el.style.background=/restaurant|cafe|fast_food|pub|bar/.test(x.category||'')?'#ed8a3c':x.category==='shop'?'#438fe4':'#df71bd';el.title=x.name;el.setAttribute('aria-label',x.name);el.textContent=x.category==='restaurant'?'♨':x.category==='hotel'?'▣':x.category==='shop'?'▢':'•';
   el.addEventListener('click',e=>{e.stopPropagation();latest.current.onChoose(x)});
   markers.push(new Marker({element:el,anchor:'bottom'}).setLngLat([x.lon,x.lat]).addTo(map));
  }
  return()=>markers.forEach(m=>m.remove());
 },[p.places,p.stops,p.selected,p.generation,p.viewportGeneration]);
 useEffect(()=>{
  const map=native.current;if(!map?.isStyleLoaded())return;
  for(const [id,coordinates,color] of [['route',p.route?.geometry?.coordinates||[],'#4285f4'],['measurement',p.measurePoints,'#008b91']] as any[]){
   const data:any={type:'FeatureCollection',features:coordinates.length>1?[{type:'Feature',properties:{},geometry:{type:'LineString',coordinates}}]:[]};
   const src:any=map.getSource(id);if(src)src.setData(data);else{map.addSource(id,{type:'geojson',data});map.addLayer({id,type:'line',source:id,layout:{'line-cap':'round','line-join':'round'},paint:{'line-color':color,'line-width':id==='route'?6:3}})}
  }
 },[p.route,p.measurePoints,p.generation]);
 return React.createElement('div',{ref:node,className:'map-canvas',role:'region','aria-label':'Interactive map'});
}
