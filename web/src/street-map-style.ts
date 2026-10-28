import {
 CenteredTextSymbolizer,CircleSymbolizer,GeomType,GroupSymbolizer,
 LineSymbolizer,OffsetTextSymbolizer,Padding,PolygonSymbolizer,linear,
 type Feature,type LabelRule,type PaintRule,
} from 'protomaps-leaflet';
import {RoadLabel} from './road-label';

// Draw our installed open-data tiles locally. Colours and hierarchy follow the
// supplied mobile reference; no remote style, font, sprite or Google data is used.
export const mapPalette={
 land:'#f5f5f5',water:'#8ed8e8',park:'#c6f0d5',wood:'#b9e7cb',
 building:'#e7e9ec',buildingEdge:'#dadde1',road:'#c9d2dc',roadEdge:'#aebbc8',
 street:'#ffffff',streetEdge:'#d5d9de',text:'#4c5865',parkText:'#31805c',
 waterText:'#527f91',rail:'#aeb4bd',boundary:'#c4c0c9',
};
const nameProps=['name:en','name'];
const font='Arial, sans-serif';
const kind=(f:Feature)=>String(f.props.kind||'');
const detail=(f:Feature)=>String(f.props.kind_detail||'');
const named=(f:Feature)=>nameProps.some(k=>typeof f.props[k]==='string'&&f.props[k]);
const parkKinds=['park','garden','golf_course','recreation_ground','village_green','playground','pitch','cemetery','national_park','protected_area'];
const woods=['wood','forest','nature_reserve'];
const grass=['grass','grassland','scrub','meadow','allotments'];
const polygon=(id:string,layer:string,fill:string,filter?:PaintRule['filter'],minzoom?:number):PaintRule=>({
 id,dataLayer:layer,minzoom,filter,symbolizer:new PolygonSymbolizer({fill}),
});
const line=(id:string,layer:string,color:string,width:number|((z:number,f?:Feature)=>number),filter:PaintRule['filter'],minzoom?:number):PaintRule=>({
 id,dataLayer:layer,minzoom,filter,symbolizer:new LineSymbolizer({color,width,lineCap:'round',lineJoin:'round'}),
});

// Separate casings and fills keep junctions continuous. Bridges are drawn after
// surface streets; tunnel colours stay subdued. Sizes are CSS pixels at each zoom.
const roadWidths={
 service:linear([[12,0],[14,1.4],[16,3],[19,6]]),
 minor:linear([[11,0],[13,1.4],[15,3.8],[17,7],[19,12]]),
 major:linear([[6,.6],[11,1.5],[13,3],[15,6.5],[17,11],[19,18]]),
 highway:linear([[3,.5],[8,1.4],[12,3],[15,7.5],[17,13],[19,20]]),
};
type RoadClass=keyof typeof roadWidths;
const roadClass=(f:Feature):RoadClass|null=>{
 if(kind(f)==='highway')return 'highway';
 if(kind(f)==='major_road')return 'major';
 if(kind(f)==='minor_road')return ['service','driveway','parking_aisle'].includes(detail(f))?'service':'minor';
 return null;
};
const roadRules=(levels:('tunnel'|'surface'|'bridge')[]):PaintRule[]=>{
 const result:PaintRule[]=[];
 for(const level of levels){
  for(const casing of [true,false])for(const cls of ['service','minor','major','highway'] as RoadClass[]){
   const main=cls==='major'||cls==='highway';
   const color=casing?mapPalette[main?'roadEdge':'streetEdge']:level==='tunnel'?'#e0e4e8':mapPalette[main?'road':'street'];
   result.push(line(`${level}-${cls}-${casing?'casing':'fill'}`,'roads',color,(z)=>{
    const w=roadWidths[cls](z);return w>0?w+(casing?(main?1.3:1):0):0;
   },(_z,f)=>roadClass(f)===cls&&(level==='bridge'?!!f.props.is_bridge:level==='tunnel'?!!f.props.is_tunnel&&!f.props.is_bridge:!f.props.is_bridge&&!f.props.is_tunnel),cls==='service'?13:cls==='minor'?12:undefined));
  }
 }
 return result;
};
const poiColor=(f?:Feature)=>{
 const k=f?kind(f):'';
 if([...parkKinds,...woods,'natural','outdoors'].includes(k))return mapPalette.parkText;
 if(['restaurant','cafe','fast_food','pub','bar','food_and_drink'].includes(k))return '#b27835';
 if(['hospital','clinic','pharmacy','doctors'].includes(k))return '#c65b66';
 if(['station','bus_stop','aerodrome','airport','ferry_terminal'].includes(k))return '#397cb5';
 if(['hotel','motel','hostel','attraction','museum','artwork','tourism','theme_park'].includes(k))return '#b45da4';
 return '#756294';
};

export function streetMapStyle(){
 const paintRules:PaintRule[]=[
  polygon('land','earth',mapPalette.land),
  {id:'landcover',dataLayer:'landcover',symbolizer:new PolygonSymbolizer({fill:(_z,f)=>{
   const k=f?kind(f):'';return k==='forest'?mapPalette.wood:['grassland','scrub'].includes(k)?'#d7f0df':k==='farmland'?'#edf0df':k==='barren'?'#f3efdf':mapPalette.land;
  }})},
  polygon('parks','landuse',mapPalette.park,(_z,f)=>parkKinds.includes(kind(f))),
  polygon('woodland','landuse',mapPalette.wood,(_z,f)=>woods.includes(kind(f))),
  polygon('grass','landuse','#d4f2df',(_z,f)=>grass.includes(kind(f))),
  polygon('commercial','landuse','#fcf3d6',(_z,f)=>['commercial','retail'].includes(kind(f)),12),
  polygon('industrial','landuse','#edeef1',(_z,f)=>['industrial','aerodrome','military','naval_base'].includes(kind(f))),
  polygon('education','landuse','#f0ede3',(_z,f)=>['school','university','college'].includes(kind(f))),
  polygon('hospital','landuse','#fae9ea',(_z,f)=>kind(f)==='hospital'),
  polygon('beach','landuse','#f6edcf',(_z,f)=>['beach','sand'].includes(kind(f))),
  polygon('pedestrian','landuse','#eeece6',(_z,f)=>['pedestrian','pier'].includes(kind(f))),
  polygon('water','water',mapPalette.water,(_z,f)=>f.geomType===GeomType.Polygon),
  line('waterways','water',mapPalette.water,linear([[9,.5],[13,1.5],[16,4],[19,10]]),(_z,f)=>f.geomType===GeomType.Line),
  {id:'buildings',dataLayer:'buildings',minzoom:14,symbolizer:new PolygonSymbolizer({fill:mapPalette.building,stroke:mapPalette.buildingEdge,width:.45})},
  line('runways','roads','#d7dce1',linear([[11,1],[15,8],[19,28]]),(_z,f)=>['runway','taxiway'].includes(detail(f)),11),
  {id:'paths',dataLayer:'roads',minzoom:13,filter:(_z,f)=>['path','other'].includes(kind(f))&&!['runway','taxiway'].includes(detail(f)),symbolizer:new LineSymbolizer({color:'#c8ceca',width:1,dash:[2,2],dashColor:'#bec9c1',dashWidth:1,lineCap:'round'})},
  ...roadRules(['tunnel','surface']),
  line('rail-bed','roads','#e2e4e7',linear([[10,.7],[15,2],[18,3]]),(_z,f)=>kind(f)==='rail'),
  {id:'railway',dataLayer:'roads',minzoom:10,filter:(_z,f)=>kind(f)==='rail',symbolizer:new LineSymbolizer({color:mapPalette.rail,dash:[4,3],dashColor:mapPalette.rail,dashWidth:linear([[10,.5],[15,1],[19,1.8]])})},
  ...roadRules(['bridge']),
  {id:'boundaries',dataLayer:'boundaries',filter:(_z,f)=>!f.props.kind_detail||Number(f.props.kind_detail)<=4,symbolizer:new LineSymbolizer({color:mapPalette.boundary,dash:[3,3],dashColor:mapPalette.boundary,dashWidth:.8})},
 ];
 const labels:LabelRule[]=[
  {id:'cities',dataLayer:'places',filter:(z,f)=>kind(f)==='locality'&&z>=Number(f.props.min_zoom||3),sort:(a,b)=>Number(a.min_zoom||0)-Number(b.min_zoom||0),symbolizer:new CenteredTextSymbolizer({labelProps:nameProps,fill:mapPalette.text,stroke:mapPalette.land,width:2,font:(z)=>`${z<10?500:400} ${z<10?13:16}px ${font}`})},
  {id:'major-road-names',dataLayer:'roads',minzoom:11,filter:(_z,f)=>['highway','major_road'].includes(kind(f))&&!!named(f),symbolizer:new RoadLabel(12,mapPalette.text,'#f5f6f7')},
  {id:'street-names',dataLayer:'roads',minzoom:14,filter:(_z,f)=>roadClass(f)==='minor'&&!!named(f),symbolizer:new RoadLabel(11,mapPalette.text,'#ffffff')},
  // The tile's POI point is the label anchor. A landuse polygon's first vertex
  // would put its name at an arbitrary corner instead of inside the park.
  {id:'park-names',dataLayer:'pois',minzoom:13,filter:(z,f)=>[...parkKinds,...woods].includes(kind(f))&&!!named(f)&&z>=Number(f.props.min_zoom||13),symbolizer:new Padding(6,new CenteredTextSymbolizer({labelProps:nameProps,fill:mapPalette.parkText,stroke:mapPalette.park,width:1.5,maxLineChars:24,lineHeight:1.2,font:`400 13px ${font}`}))},
  {id:'poi-names',dataLayer:'pois',minzoom:14,filter:(z,f)=>!!named(f)&&![...parkKinds,...woods,'administrative','platform','entrance','address'].includes(kind(f))&&z>=Math.max(14,Number(f.props.min_zoom||15)+1),sort:(a,b)=>Number(a.min_zoom||15)-Number(b.min_zoom||15),symbolizer:new Padding(10,new GroupSymbolizer([
   new CircleSymbolizer({radius:3.5,fill:(_z,f)=>poiColor(f),stroke:'#ffffff',width:1.5}),
   new OffsetTextSymbolizer({labelProps:nameProps,fill:(_z,f)=>poiColor(f),stroke:'#ffffff',width:2,maxLineChars:24,lineHeight:1.2,font:`400 12px ${font}`,offsetX:8,offsetY:0}),
  ]))},
  {id:'neighbourhoods',dataLayer:'places',minzoom:13,filter:(_z,f)=>['neighbourhood','suburb','quarter','borough'].includes(kind(f)),symbolizer:new CenteredTextSymbolizer({labelProps:nameProps,fill:'#6e7881',stroke:mapPalette.land,width:2,font:`400 12px ${font}`})},
  {id:'water-names',dataLayer:'water',filter:(_z,f)=>f.geomType===GeomType.Point&&!!named(f),symbolizer:new CenteredTextSymbolizer({labelProps:nameProps,fill:mapPalette.waterText,font:`italic 12px ${font}`})},
  {id:'countries',dataLayer:'places',maxzoom:8,filter:(_z,f)=>kind(f)==='country',symbolizer:new CenteredTextSymbolizer({labelProps:nameProps,fill:'#687582',stroke:mapPalette.land,width:2,font:`500 13px ${font}`})},
  {id:'regions',dataLayer:'places',minzoom:5,maxzoom:10,filter:(_z,f)=>kind(f)==='region',symbolizer:new CenteredTextSymbolizer({labelProps:nameProps,fill:'#85909b',stroke:mapPalette.land,width:1.5,font:`400 11px ${font}`})},
 ];
 return {backgroundColor:mapPalette.land,paintRules,labelRules:labels};
}
