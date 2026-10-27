import React from 'react';
import {RasterMap} from './raster-map';
export type Place={id?:string,key?:string,name:string,address:string,lat:number,lon:number,list?:string,notes?:string,category?:string,hours?:string,phone?:string,website?:string,osm?:string};
export const distance=(a:number[],b:number[])=>{const rad=Math.PI/180,lat=(b[1]-a[1])*rad,lon=(b[0]-a[0])*rad;return 12742000*Math.asin(Math.sqrt(Math.sin(lat/2)**2+Math.cos(a[1]*rad)*Math.cos(b[1]*rad)*Math.sin(lon/2)**2))};
export const fmtDistance=(m:number)=>m<1000?Math.round(m)+' m':(m/1000).toFixed(1)+' km';
export function MapView(p:any){return React.createElement(RasterMap,p)}
