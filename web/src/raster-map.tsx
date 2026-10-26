import React,{useEffect,useRef} from 'react';
import * as L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import type {MapController} from './map-controller';
import type {Place} from './map';

export function RasterMap(p:any){
 const node=useRef<HTMLDivElement>(null),native=useRef<L.Map|null>(null),latest=useRef(p),previousSatellite=useRef(false),vectorZoom=useRef(12);latest.current=p;
 useEffect(()=>{
  const saved=p.camera.current,center=saved?.center||(p.selected?[p.selected.lon,p.selected.lat]:[-.09,51.505]);
  const map=L.map(node.current!,{zoomControl:false,attributionControl:false,preferCanvas:false}).setView([center[1],center[0]],saved?.zoom||(p.selected?15:12));native.current=map;
  // SVG/DOM tiles and overlays work without WebGL, including disabled GPU drivers.
  L.control.attribution({position:'bottomleft',prefix:false}).addTo(map);L.control.zoom({position:'bottomright'}).addTo(map);
  const controller:MapController={renderer:'raster',native:map,getCenter:()=>{const c=map.getCenter();return {lat:c.lat,lng:c.lng,toArray:()=>[c.lng,c.lat]}},getZoom:()=>map.getZoom(),getBearing:()=>0,getPitch:()=>0,
   easeTo:o=>{const c=o.center;if(c)map.setView([c[1],c[0]],o.zoom??map.getZoom(),{animate:false});else if(o.zoom!==undefined)map.setZoom(o.zoom,{animate:false});if(c&&o.padding){const d=o.padding;map.panBy([(d.right-d.left)/2,(d.bottom-d.top)/2],{animate:false})}},
   fitBounds:(b,o)=>map.fitBounds([[b[0][1],b[0][0]],[b[1][1],b[1][0]]],{maxZoom:o.maxZoom,paddingTopLeft:[o.padding.left,o.padding.top],paddingBottomRight:[o.padding.right,o.padding.bottom],animate:false}),
   loaded:()=>true,once:(e,f)=>{map.once(e==='idle'?'moveend':e,f)},off:(e,f)=>{map.off(e==='idle'?'moveend':e,f)},getStyle:()=>undefined,querySourceFeatures:()=>[]};
  p.mapRef.current=controller;
  map.on('click',(e:L.LeafletMouseEvent)=>latest.current.onPoint({lat:e.latlng.lat,lon:e.latlng.lng}));
  map.on('moveend',()=>latest.current.onMove());
  const resize=new ResizeObserver(()=>map.invalidateSize());resize.observe(node.current!);latest.current.onReady();
  return()=>{resize.disconnect();latest.current.onMove();map.remove();native.current=null;if(p.mapRef.current===controller)p.mapRef.current=null};
 },[]);
 useEffect(()=>{
  const map=native.current;if(!map)return;let reported=false;
  const satellite=p.mapType==='satellite';
  const url=satellite?'https://gibs.earthdata.nasa.gov/wmts/epsg3857/best/BlueMarble_ShadedRelief/default/GoogleMapsCompatible_Level8/{z}/{y}/{x}.jpeg':p.config.tiles;
  // Escape only after Leaflet expands z/x/y, preserving configured tile URLs.
  const layer=L.tileLayer('',{tileSize:256,maxZoom:satellite?8:19,attribution:satellite?'Imagery: NASA Blue Marble (500 m)':'',detectRetina:false});
  layer.getTileUrl=(coords:L.Coords)=>p.base+'/api/map-resource?url='+encodeURIComponent(L.Util.template(url,{z:coords.z,x:coords.x,y:coords.y}));
  const attribution=document.createElement('span');attribution.textContent=satellite?'Imagery: NASA Blue Marble (500 m)':p.config.attribution;
  const credit=!satellite&&!p.config.customTiles?'© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap contributors</a>':attribution.innerHTML;map.attributionControl?.addAttribution(credit);
  layer.on('tileerror',()=>{if(!reported){reported=true;latest.current.onError('Map tiles could not load. Check your connection, then move or reload the map.')}});
  layer.on('load',()=>latest.current.onReady());layer.addTo(map);
  if(satellite){vectorZoom.current=map.getZoom();map.setZoom(Math.min(8,map.getZoom()))}else if(previousSatellite.current)map.setZoom(vectorZoom.current);previousSatellite.current=satellite;
  return()=>{layer.remove();map.attributionControl?.removeAttribution(credit)};
 },[p.mapType]);
 useEffect(()=>{
  const map=native.current;if(!map)return;const group=L.layerGroup().addTo(map),unique=new Map<string,Place>();
  for(const x of [...(p.selected?[p.selected]:[]),...p.stops,...p.places])if(!unique.has(x.lat+','+x.lon))unique.set(x.lat+','+x.lon,x);
  const occupied:L.Point[]=[];for(const x of unique.values()){
   const point=map.latLngToContainerPoint([x.lat,x.lon]);if(x!==p.selected&&occupied.some(v=>v.distanceTo(point)<32))continue;occupied.push(point);
   const el=document.createElement('button');el.type='button';el.className='place-pin'+(x===p.selected?' place-pin--selected':'');el.title=x.name;el.setAttribute('aria-label',x.name);el.textContent='•';el.style.background=/restaurant|cafe|fast_food|pub|bar/.test(x.category||'')?'#ed8a3c':x.category==='shop'?'#438fe4':'#df71bd';
   el.addEventListener('click',e=>{L.DomEvent.stopPropagation(e);latest.current.onChoose(x)});
   L.marker([x.lat,x.lon],{icon:L.divIcon({html:el,className:'raster-place-marker',iconSize:[32,40],iconAnchor:[16,40]})}).addTo(group);
  }
  return()=>{group.remove()};
 },[p.places,p.stops,p.selected,p.generation,p.viewportGeneration]);
 useEffect(()=>{
  const map=native.current;if(!map)return;const group=L.layerGroup().addTo(map);
  for(const [coordinates,color,weight] of [[p.route?.geometry?.coordinates||[],'#4285f4',6],[p.measurePoints,'#008b91',3]] as any[])if(coordinates.length>1)L.polyline(coordinates.map((c:number[])=>[c[1],c[0]]),{color,weight}).addTo(group);
  return()=>{group.remove()};
 },[p.route,p.measurePoints]);
 return React.createElement('div',{ref:node,className:'map-canvas raster-map',role:'region','aria-label':'Interactive map'});
}
