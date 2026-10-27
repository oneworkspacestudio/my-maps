import * as L from 'leaflet';

// The product controls work with either renderer; graphics support must not
// determine whether search, routes, private pins or GPS can be used.
export interface MapController {
 renderer:'vector'|'raster'; native:L.Map;
 getCenter():{lat:number,lng:number,toArray():[number,number]};
 getZoom():number; getBearing():number; getPitch():number;
 easeTo(options:any):void; fitBounds(bounds:any,options:any):void;
 loaded():boolean; once(event:'idle',fn:()=>void):void; off(event:'idle',fn:()=>void):void;
 getStyle():any; querySourceFeatures(source:string,options:any):any[];
}
export interface LocationMarker {remove():void;setLngLat(ll:[number,number]):LocationMarker;addTo(map:MapController):LocationMarker}
export function locationMarker(map:MapController,ll:[number,number]):LocationMarker {
 const dot=document.createElement('div');dot.className='location-dot';
 const marker=L.marker([ll[1],ll[0]],{icon:L.divIcon({html:dot,className:'gps-marker',iconSize:[20,20],iconAnchor:[10,10]}),interactive:false});
 const wrapper:LocationMarker={remove:()=>{marker.remove()},setLngLat:xy=>{marker.setLatLng([xy[1],xy[0]]);return wrapper},addTo:m=>{marker.addTo(m.native as L.Map);return wrapper}};return wrapper;
}
