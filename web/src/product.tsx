import React,{useEffect,useRef,useState} from 'react';

import {locationMarker as makeLocationMarker,type LocationMarker,type MapController} from './map-controller';

import {Maps} from './generated/Maps';

import {MapInstall} from './map-install';

import {DirectionsInstall} from './directions-install';
import {MapView,distance,fmtDistance,type Place} from './map';

const duration=(n:number)=>{const minutes=Math.max(1,Math.round(n/60));return minutes>=60?Math.floor(minutes/60)+' hr'+(minutes%60?' '+minutes%60+' min':''):minutes+' min'};

const turnIcon=(step:any)=>step?.maneuver?.type==='arrive'?'Flag':step?.maneuver?.type==='depart'?'Navigation':step?.maneuver?.type==='roundabout'?'RotateCcw':step?.maneuver?.modifier?.includes('left')?'CornerUpLeft':step?.maneuver?.modifier?.includes('right')?'CornerUpRight':'ArrowUp';
const defaultVehicle=(mode:string)=>mode==='bike'?'touring':mode==='foot'?'walking':'car';

const point=(x:any,i:number):Place=>({key:String(i),name:x.properties.name||x.properties.postcode||x.properties.street||x.properties.city||'Place',address:[(x.properties.housenumber?[x.properties.housenumber,x.properties.street].filter(Boolean).join(' '):x.properties.street),x.properties.city,x.properties.state,x.properties.country].filter(Boolean).join(', '),lon:x.geometry.coordinates[0],lat:x.geometry.coordinates[1],category:x.properties.osm_value||'Place'});

const mapFeatures=(features:any):Place[]=>(Array.isArray(features)?features:[]).filter(x=>x?.geometry?.type==='Point'&&Array.isArray(x.geometry.coordinates)&&x.geometry.coordinates.length===2&&x.geometry.coordinates.every(Number.isFinite)&&Math.abs(x.geometry.coordinates[0])<=180&&Math.abs(x.geometry.coordinates[1])<=90&&x.properties).map(point);

const categories=[{key:'restaurants',title:'Restaurants',icon:'Utensils'},{key:'shopping',title:'Shopping',icon:'ShoppingBag'},{key:'hotels',title:'Hotels',icon:'BedDouble'},{key:'attractions',title:'Things to do',icon:'Ticket'},{key:'fuel',title:'Petrol',icon:'Fuel'},{key:'parking',title:'Parking',icon:'CircleParking'}];

const safeURL=(raw?:string)=>{try{const u=new URL(raw||'');return ['http:','https:'].includes(u.protocol)&&!u.username&&!u.password?u.href:''}catch{return ''}};

const categoryName=(p:Place)=>p.category?.replaceAll('_',' ')||'Place';

const iconFor=(p:Place)=>/restaurant|cafe|pub|bar|fast_food/.test(p.category||'')?'Utensils':/hotel|hostel|guest_house/.test(p.category||'')?'BedDouble':/shop/.test(p.category||'')?'ShoppingBag':'MapPin';

export function Product(p:any){

 const [config,setConfig]=useState<any>(null),[tab,setTab]=useState('explore'),[results,setResults]=useState<Place[]>([]),[places,setPlaces]=useState<Place[]>([]),[selected,setSelected]=useState<Place|null>(null),[stops,setStops]=useState<Place[]>([]),[routes,setRoutes]=useState<any[]>([]),[routeIndex,setRouteIndex]=useState(0),[routeResults,setRouteResults]=useState<Place[]>([]),[searchTarget,setSearchTarget]=useState('end'),[hasStart,setHasStart]=useState(false);

 const [searchOpen,setSearchOpen]=useState(false),[recents,setRecents]=useState<any[]>([]),[suggestions,setSuggestions]=useState<Place[]>([]),[searchStatus,setSearchStatus]=useState(''),[suggestionPending,setSuggestionPending]=useState(false);

 const [saveTarget,setSaveTarget]=useState<Place|null>(null),[shortcutIntent,setShortcutIntent]=useState('');

 const [routeBusy,setRouteBusy]=useState(false),[optionsOpen,setOptionsOpen]=useState(false),[needsCoverage,setNeedsCoverage]=useState(false),[navProgress,setNavProgress]=useState<any>(null);

 const routeEpoch=useRef(0),routeAbort=useRef<AbortController|null>(null),routeCurrent=useRef<any>(null),navIndex=useRef(0),navStopIndex=useRef(1),routeStops=useRef<Place[]>([]);

 const cancelRoute=()=>{routeEpoch.current++;routeAbort.current?.abort();routeAbort.current=null;setRouteBusy(false)};

 const ignoreSearchPop=useRef(false);

 const suggestionGeneration=useRef(0),searchHistory=useRef(false),historyQueue=useRef<Promise<any>>(Promise.resolve()),recentGeneration=useRef(0);

 const [resultTitle,setResultTitle]=useState(''),[attempted,setAttempted]=useState(false),[emptyMessage,setEmptyMessage]=useState('');

 const [expanded,setExpanded]=useState(false),[overlay,setOverlay]=useState(''),[layerMessage,setLayerMessage]=useState(''),[category,setCategory]=useState(''),[generation,setGeneration]=useState(0),[moved,setMoved]=useState(0),[measuring,setMeasuring]=useState(false),[measurePoints,setMeasurePoints]=useState<number[][]>([]),[navigating,setNavigating]=useState(false);

 const camera=useRef<any>(null),mapRef=useRef<MapController|null>(null),watch=useRef<number|null>(null),locationMarker=useRef<LocationMarker|null>(null),markerMap=useRef<MapController|null>(null),alive=useRef(true),selection=useRef(0),searchGeneration=useRef(0),f=p.fields;

 const center=()=>{const c=mapRef.current?.getCenter();return {lat:c?.lat??51.505,lon:c?.lng??-.09}};

 const coords=()=>{const c=center();return 'lat='+c.lat+'&lon='+c.lon};

 const close=()=>{setOverlay('');setSaveTarget(null);setLayerMessage('')};

 const leaveSearch=()=>{setSearchOpen(false);if(searchHistory.current){searchHistory.current=false;ignoreSearchPop.current=true;history.back()}};

 const stopFollowing=()=>{if(watch.current!==null)navigator.geolocation.clearWatch(watch.current);watch.current=null;setNavigating(false);setNavProgress(null);navIndex.current=0;navStopIndex.current=1};

 useEffect(()=>{const pop=()=>{searchHistory.current=false;setSearchOpen(false);if(ignoreSearchPop.current){ignoreSearchPop.current=false;return}setShortcutIntent('');setSaveTarget(null);setOverlay('')};window.addEventListener('popstate',pop);return()=>window.removeEventListener('popstate',pop)},[]);

 const loadRecents=async()=>{const stamp=++recentGeneration.current;const items=await p.request('/api/recents');if(alive.current&&stamp===recentGeneration.current)setRecents(items)};

 const remember=(value:any)=>{const operation=historyQueue.current.then(async()=>{try{await p.request('/api/recents',{method:'POST',body:JSON.stringify(value)});await loadRecents()}catch{if(alive.current)p.setNotice('Recent search could not be saved. Your map is still available.')}});historyQueue.current=operation;return operation};

 const loadPlaces=async()=>{const value=await p.request('/api/places');if(alive.current)setPlaces(value);return value};

 useEffect(()=>{alive.current=true;p.request('/api/config').then((v:any)=>alive.current&&setConfig(v)).catch((e:any)=>p.setError(e.message));loadPlaces().catch((e:any)=>p.setError(e.message));loadRecents().catch(()=>{});return()=>{alive.current=false;routeAbort.current?.abort();routeEpoch.current++;selection.current++;searchGeneration.current++;if(watch.current!==null)navigator.geolocation.clearWatch(watch.current);locationMarker.current?.remove()}},[]);

 useEffect(()=>{const q=new URLSearchParams(location.hash.slice(1));if(q.has('route')){try{const raw=q.get('route')||'';if(raw.length>6000)throw Error();const shared=JSON.parse(raw);if(!['car','bike','foot'].includes(shared.mode)||!Array.isArray(shared.stops)||shared.stops.length<2||shared.stops.length>8)throw Error();const list=shared.stops.map((v:any)=>{if(!Number.isFinite(v.lat)||!Number.isFinite(v.lon)||Math.abs(v.lat)>85||Math.abs(v.lon)>180||typeof v.name!=='string'||v.name.length>100)throw Error();return {lat:v.lat,lon:v.lon,name:v.name,address:''}});setStops(list);setHasStart(true);p.setField('mode',shared.mode);p.setField('vehicle',defaultVehicle(shared.mode));p.setField('startQuery',list[0].name);p.setField('endQuery',list.at(-1).name);setTab('directions');setExpanded(true)}catch{p.setError('This shared route link is invalid.')}return}const lat=Number(q.get('lat')),lon=Number(q.get('lon'));if(q.has('lat')&&q.has('lon')&&Number.isFinite(lat)&&Math.abs(lat)<=90&&Number.isFinite(lon)&&Math.abs(lon)<=180)choose({name:q.get('name')||'Shared place',address:'',lat,lon})},[]);

 useEffect(()=>{routeCurrent.current=routes[routeIndex]||null;const r=routes[routeIndex],map=mapRef.current;if(!map||!r)return;const xy=r.geometry.coordinates;let west=180,south=90,east=-180,north=-90;for(const c of xy){west=Math.min(west,c[0]);east=Math.max(east,c[0]);south=Math.min(south,c[1]);north=Math.max(north,c[1])}const mobile=innerWidth<=700;map.fitBounds([[west,south],[east,north]],{padding:mobile?{top:160,bottom:270,left:35,right:70}:{top:140,bottom:70,left:470,right:90},maxZoom:16});},[routes,routeIndex]);

 useEffect(()=>{if(p.route==='#settings'){cancelRoute();stopFollowing();close()}},[p.route]);

 useEffect(()=>{document.querySelector('.stack--sheet-content')?.scrollTo({top:0})},[tab,category,selected?.name]);

 useEffect(()=>{

  if(!searchOpen)return;

  const stamp=++suggestionGeneration.current,query=String(f.query||'').trim(),controller=new AbortController();

  setSuggestions([]);setSearchStatus('');setSuggestionPending(query.length>=2);

  if(query.length<2)return;

  const timer=setTimeout(()=>{p.request('/api/search?q='+encodeURIComponent(query)+'&'+coords(),{signal:controller.signal}).then((r:any)=>{if(alive.current&&stamp===suggestionGeneration.current){const found=mapFeatures(r.features);setSuggestions(found);setSearchStatus(found.length?'':'No matching places. Try a town, address or postcode.')}}).catch((e:any)=>{if(e.name!=='AbortError'&&alive.current&&stamp===suggestionGeneration.current)setSearchStatus('Place suggestions are unavailable. You can retry your search.')}).finally(()=>{if(alive.current&&stamp===suggestionGeneration.current)setSuggestionPending(false)})},400);

  return()=>{clearTimeout(timer);controller.abort();suggestionGeneration.current++};

 },[searchOpen,f.query]);

 useEffect(()=>{const keyboard=()=>{document.documentElement.style.setProperty('--app-height',(window.visualViewport?.height||innerHeight)+'px')};keyboard();window.visualViewport?.addEventListener('resize',keyboard);window.addEventListener('resize',keyboard);return()=>{window.visualViewport?.removeEventListener('resize',keyboard);window.removeEventListener('resize',keyboard);document.documentElement.style.removeProperty('--app-height')}},[]);

 useEffect(()=>{const back=(event:KeyboardEvent)=>{if(event.key==='Escape'&&searchOpen&&overlay!=='save'){event.preventDefault();leaveSearch();blurSearch()}};window.addEventListener('keydown',back);return()=>window.removeEventListener('keydown',back)},[searchOpen,overlay]);

 const samePlace=(a:Place,b:Place)=>a.id&&b.id?a.id===b.id:a.name===b.name&&distance([a.lon,a.lat],[b.lon,b.lat])<5;

 const savedVersion=(value:Place)=>{const saved=places.find(x=>samePlace(x,value));return saved?{...value,...saved}:value};

 const listNames=[...new Set(['Favourites','Home','Work',...places.map(x=>x.list||'Favourites')])].filter((x,i,all)=>all.findIndex(v=>v.toLocaleLowerCase()===x.toLocaleLowerCase())===i);

 const openSave=(value:Place,list?:string)=>{const x=savedVersion(value);blurSearch();setSaveTarget(x);p.setField('saveList',list||x.list||'Favourites');p.setField('newList','');p.setField('saveNotes',selected&&samePlace(selected,x)?f.notes??x.notes??'':x.notes||'');setOverlay('save')};

 const openSearch=()=>{if(!searchOpen){history.pushState({mapsSearch:true},'',location.href);searchHistory.current=true}setSearchOpen(true);close()};

 const choose=(value:Place,collapse=true,offerShortcut=true)=>{cancelRoute();const saved=places.find(place=>place.name===value.name&&distance([place.lon,place.lat],[value.lon,value.lat])<5);const x=saved?{...value,...saved}:value;stopFollowing();setRoutes([]);selection.current++;leaveSearch();setTab('explore');setSelected(x);p.setField('list',x.list||'Favourites');p.setField('notes',x.notes||'');if(collapse)setExpanded(false);if(offerShortcut&&shortcutIntent){setShortcutIntent('');openSave(x,shortcutIntent)}mapRef.current?.easeTo({center:[x.lon,x.lat],zoom:Math.max(mapRef.current.getZoom(),15),padding:innerWidth<=700?{top:140,bottom:230,left:40,right:70}:{top:60,bottom:30,left:440,right:30}})};

 const pick=async(x:any)=>{if(measuring){setMeasurePoints(v=>[...v,[x.lon,x.lat]]);return}choose({lat:x.lat,lon:x.lon,name:x.name||'Dropped pin',address:''});if(x.name)return;const stamp=selection.current;try{const r=await p.request('/api/reverse?lat='+x.lat+'&lon='+x.lon);if(alive.current&&stamp===selection.current&&r.features?.length)choose({...point(r.features[0],0),lat:x.lat,lon:x.lon},false,false)}catch(e:any){if(alive.current&&stamp===selection.current)p.setError(e.message)}};

 const setMainTab=(value:string)=>{cancelRoute();stopFollowing();leaveSearch();selection.current++;searchGeneration.current++;setTab(value);setSelected(null);setExpanded(value!=='explore');close()};

 const blurSearch=()=>{if(document.activeElement instanceof HTMLInputElement)document.activeElement.blur()};

 const searchPlaces=async(query?:string)=>{

  const q=String(query??f.query??'').trim();if(q.length<2)throw Error('Enter a place or address');

  const categoryQuery:Record<string,string>={restaurants:'restaurants',restaurant:'restaurants',shopping:'shopping',hotels:'hotels',hotel:'hotels',attractions:'attractions',petrol:'fuel',parking:'parking'};

  leaveSearch();blurSearch();void remember({query:q,name:'',address:''});

  if(categoryQuery[q.toLowerCase()]){await nearby(categoryQuery[q.toLowerCase()]);return}

  const stamp=++searchGeneration.current;selection.current++;setSelected(null);setResults([]);setCategory('');setResultTitle(q);setAttempted(true);setEmptyMessage('Finding places...');setTab('explore');setExpanded(false);p.setBusy('Finding places...');

  try{const response=await p.request('/api/search?q='+encodeURIComponent(q)+'&'+coords());if(!alive.current||stamp!==searchGeneration.current)return;const found=mapFeatures(response.features);setResults(found);setEmptyMessage('No places found. Try a town, address or postcode.');if(found.length){choose(found[0],false);void remember({query:'',name:found[0].name,address:found[0].address,lat:found[0].lat,lon:found[0].lon})}}

  catch(e){if(!alive.current||stamp!==searchGeneration.current)return;setEmptyMessage('Search is temporarily unavailable. Try again shortly.');throw e}

 };

 const nearby=async(key:string)=>{stopFollowing();setRoutes([]);leaveSearch();const stamp=++searchGeneration.current;selection.current++;setSelected(null);setCategory(key);setResultTitle('');setTab('explore');setResults([]);setAttempted(true);setEmptyMessage('Finding nearby places…');setExpanded(false);blurSearch();p.setBusy('Finding nearby places…');const c=center();if(mapRef.current?.renderer==='vector'&&mapRef.current.getZoom()<14.5)mapRef.current.easeTo({zoom:15});let response:any;try{response=await p.request('/api/nearby?'+coords()+'&category='+key)}catch(e){setEmptyMessage('The local place index is unavailable. Retry or check Map downloads.');throw e};if(!alive.current||stamp!==searchGeneration.current)return;const unique=new Map<string,Place>();for(const x of response.elements||[]){const tags=x.tags||{},lat=x.lat??x.center?.lat,lon=x.lon??x.center?.lon;if(!Number.isFinite(lat)||!Number.isFinite(lon)||Math.abs(lat)>90||Math.abs(lon)>180)continue;const place:Place={key:x.type+'/'+x.id,name:tags.name||tags.brand||tags.amenity?.replaceAll('_',' ')||key,address:[tags['addr:housenumber'],tags['addr:street'],tags['addr:city']].filter(Boolean).join(', '),lat,lon,category:tags.amenity||tags.tourism||(tags.shop?'shop':tags.public_transport)||key,hours:tags.opening_hours,phone:tags.phone||tags['contact:phone'],website:safeURL(tags.website||tags['contact:website']),osm:['node','way','relation'].includes(x.type)?'https://www.openstreetmap.org/'+x.type+'/'+x.id:''};if(![...unique.values()].some(v=>v.name===place.name&&distance([v.lon,v.lat],[lon,lat])<35))unique.set(place.key!,place)}const list=[...unique.values()].sort((a,b)=>distance([c.lon,c.lat],[a.lon,a.lat])-distance([c.lon,c.lat],[b.lon,b.lat]));setResults(list);if(list.length){const xs=list.map(x=>x.lon),ys=list.map(x=>x.lat);mapRef.current?.fitBounds([[Math.min(...xs),Math.min(...ys)],[Math.max(...xs),Math.max(...ys)]],{padding:innerWidth<=700?{top:185,bottom:160,left:35,right:85}:{top:155,bottom:90,left:470,right:85},maxZoom:16})}setEmptyMessage('No mapped places found nearby. Move the map or choose another category.')};

 const updatePosition=(pos:GeolocationPosition)=>{

  if(!alive.current||!mapRef.current)return;

  const ll:[number,number]=[pos.coords.longitude,pos.coords.latitude];

  if(!locationMarker.current)locationMarker.current=makeLocationMarker(mapRef.current,ll).addTo(mapRef.current);else locationMarker.current.setLngLat(ll).addTo(mapRef.current);

  mapRef.current.easeTo({center:ll,zoom:17});

  const route=routeCurrent.current;if(!route)return;while(navStopIndex.current<routeStops.current.length-1&&distance(ll,[routeStops.current[navStopIndex.current].lon,routeStops.current[navStopIndex.current].lat])<50)navStopIndex.current++;

  const xy=route.geometry.coordinates;let best=Infinity,nearest=navIndex.current;

  // Project onto each remaining segment, rather than jumping between sparse graph nodes.

  for(let i=Math.max(0,navIndex.current-3);i<xy.length-1;i++){const a=xy[i],b=xy[i+1],scale=Math.cos(ll[1]*Math.PI/180),dx=(b[0]-a[0])*scale,dy=b[1]-a[1],d=dx*dx+dy*dy,t=d?Math.max(0,Math.min(1,((ll[0]-a[0])*scale*dx+(ll[1]-a[1])*dy)/d)):0,projected=[a[0]+(b[0]-a[0])*t,a[1]+(b[1]-a[1])*t],away=distance(ll,projected);if(away<best){best=away;nearest=i}}

  if(pos.coords.accuracy>100){setNavProgress({instruction:'Waiting for a more accurate location',remaining:'GPS accuracy '+Math.round(pos.coords.accuracy)+' m',offRoute:false});return}

  navIndex.current=Math.max(navIndex.current,nearest);let remaining=distance(ll,xy[nearest+1]);for(let i=nearest+1;i<xy.length-1;i++)remaining+=distance(xy[i],xy[i+1]);

  const flat=route.legs.flatMap((l:any)=>l.steps),next=flat.find((step:any)=>Number.isFinite(step.index)&&step.index>nearest)||flat.at(-1);

  let toTurn=next?distance(ll,xy[nearest+1]):0;for(let i=nearest+1;next&&i<next.index;i++)toTurn+=distance(xy[i],xy[i+1]);

  const arrived=remaining<35&&distance(ll,xy.at(-1))<45;

  setNavProgress({icon:arrived?'Flag':turnIcon(next),instruction:arrived?'You have arrived':best>Math.max(80,pos.coords.accuracy*2)?'Off route — update directions':(next?.instruction||'Continue on the route')+' in '+fmtDistance(toTurn),remaining:arrived?'Destination reached':fmtDistance(remaining)+' remaining · '+duration(route.duration*Math.min(1,remaining/route.distance)),offRoute:best>Math.max(80,pos.coords.accuracy*2),location:{name:'Your location',address:'',lon:ll[0],lat:ll[1]}});

  if(arrived){if(watch.current!==null)navigator.geolocation.clearWatch(watch.current);watch.current=null}

 };

 const currentLocation=async()=>{

  if(!navigator.geolocation)throw Error('Location is unavailable. Choose a starting point in Directions.');

  const position=await new Promise<GeolocationPosition>((resolve,reject)=>navigator.geolocation.getCurrentPosition(resolve,()=>reject(Error('Location could not be used. Choose a starting point in Directions.')),{enableHighAccuracy:true,timeout:15000}));

  return {name:'Your location',address:'',lat:position.coords.latitude,lon:position.coords.longitude} as Place;

 };

 const locate=(continuous=false)=>{if(!navigator.geolocation)throw Error('Location is unavailable in this browser');const fail=(e:GeolocationPositionError)=>{if(!alive.current)return;stopFollowing();p.setError(e.code===1?'Location permission was declined. You can choose a starting point.':'Could not locate this device. Choose a starting point on the map.')};if(continuous){stopFollowing();watch.current=navigator.geolocation.watchPosition(updatePosition,fail,{enableHighAccuracy:true,maximumAge:3000,timeout:20000});setNavigating(true);setExpanded(false)}else navigator.geolocation.getCurrentPosition(pos=>{updatePosition(pos);setExpanded(false)},fail,{enableHighAccuracy:true,timeout:15000})};

 const calculateRoute=async(list:Place[],overrides:any={})=>{

  cancelRoute();stopFollowing();setNeedsCoverage(false);const stamp=routeEpoch.current,controller=new AbortController();routeAbort.current=controller;setRouteBusy(true);setRoutes([]);

  const mode=overrides.mode||f.mode||'car';

  try{const value=await p.request('/api/route',{method:'POST',signal:controller.signal,body:JSON.stringify({mode,vehicle:overrides.vehicle||f.vehicle||defaultVehicle(mode),preference:f.preference||'fastest',avoidTolls:!!f.avoidTolls,avoidMotorways:!!f.avoidMotorways,avoidUnpaved:!!f.avoidUnpaved,alternatives:f.alternatives!==false,points:list.map(x=>[x.lon,x.lat])})});

   if(!alive.current||stamp!==routeEpoch.current)return false;

   if(value.code!=='Ok'||!value.routes?.length)throw Error('No route connects these stops. Try another travel mode or move a stop onto a road.');

   routeCurrent.current=value.routes[0];routeStops.current=list;setRoutes(value.routes);setRouteIndex(0);setSelected(null);setExpanded(false);return true;

  }catch(e:any){if(e.name==='AbortError'||!alive.current||stamp!==routeEpoch.current)return false;if(e.code==='directions_download'||e.code==='directions_coverage'){setNeedsCoverage(true);setExpanded(true)}throw e}finally{if(alive.current&&stamp===routeEpoch.current){setRouteBusy(false);routeAbort.current=null}}

 };

 const routeFromLocation=async(destination:Place)=>{

  cancelRoute();const stamp=routeEpoch.current;setTab('directions');setSelected(null);setStops([destination]);setHasStart(false);p.setField('endQuery',destination.name);setExpanded(true);leaveSearch();

  if(!config?.routingAvailable){setNeedsCoverage(true);return}

  p.setBusy('Finding your location...');const start=await currentLocation();if(!alive.current||stamp!==routeEpoch.current)return;

  const list=[start,destination];setStops(list);setHasStart(true);p.setField('startQuery',start.name);if(await calculateRoute(list))locate(true);

 };

 const action=async(a:string)=>{const id=p.selected.current;if(['apps','closeApps','logout','theme','backup','member'].includes(a)){if(a!=='theme')close();return p.common(a)}p.setError('');p.setNotice('');try{switch(a){

  case'openSearch':setShortcutIntent('');openSearch();break;

  case'closeSearch':setShortcutIntent('');leaveSearch();blurSearch();break;

  case'mapBack':location.hash='';close();break;

  case'clearQuery':p.setField('query','');document.querySelector<HTMLInputElement>('.stack--search-view input')?.focus();break;

  case'clearSearch':cancelRoute();stopFollowing();setStops([]);setHasStart(false);p.setField('startQuery','');p.setField('endQuery','');selection.current++;searchGeneration.current++;setSelected(null);setResults([]);setAttempted(false);setCategory('');setResultTitle('');p.setField('query','');setTab('explore');setExpanded(false);setRoutes([]);break;

  case'clearRecents':await historyQueue.current;await p.request('/api/recents',{method:'DELETE'});await loadRecents();break;

  case'recentSearch':{const item=recents[Number(id)];if(item){if(Number.isFinite(item.lat)&&Number.isFinite(item.lon)){choose(item);void remember(item)}else{p.setField('query',item.query);await searchPlaces(item.query)}}break}

  case'suggestion':{const item=suggestions.find(x=>x.key===id);if(item){searchGeneration.current++;setResults([item]);setAttempted(true);setCategory('');setResultTitle(item.name);choose(item);blurSearch();void remember({query:'',name:item.name,address:item.address,lat:item.lat,lon:item.lon})}break}

  case'quickStart':if(selected)await routeFromLocation(selected);break;

  case'explore':case'saved':case'contribute':setMainTab(a);break;

  case'directions':setMainTab('directions');break;

  case'profile':setOverlay('profile');break;case'mapInstall':case'layers':cancelRoute();stopFollowing();setOverlay('install');break;case'coordinates':setOverlay('coordinates');break;case'closeOverlay':if(overlay!=='save'||!p.busy)close();break;

  case'expand':setExpanded(v=>!v);break;

  case'compass':mapRef.current?.easeTo({bearing:0,pitch:0});break;

  case'search':await searchPlaces();break;

  case'goCoordinates':{const pair=String(f.coordinates||'').split(',').map(Number);if(pair.length!==2||!String(f.coordinates).split(',').every((x:string)=>x.trim())||!Number.isFinite(pair[0])||Math.abs(pair[0])>90||!Number.isFinite(pair[1])||Math.abs(pair[1])>180)throw Error('Use latitude, longitude within valid ranges');close();choose({name:'Dropped pin',lat:pair[0],lon:pair[1],address:''});break}

  case'category':await nearby(id);break;case'restaurants':await nearby('restaurants');break;case'refresh':if(category)await nearby(category);else await searchPlaces();break;

  case'place':{const x=results.find(x=>x.key===id);if(x){choose(x);void remember({query:'',name:x.name,address:x.address,lat:x.lat,lon:x.lon})}break}case'savedPlace':{const x=places.find(x=>x.id===id);if(x){choose(x);void remember({query:'',name:x.name,address:x.address,lat:x.lat,lon:x.lon})}break}

  case'deselect':selection.current++;setSelected(null);setExpanded(false);break;

  case'home':case'work':{const label=a==='home'?'Home':'Work',place=places.find(x=>x.list?.trim().toLowerCase()===a);if(place){setShortcutIntent('');choose(place);blurSearch()}else if(selected){setShortcutIntent('');openSave(selected,label)}else{setShortcutIntent(label);p.setField('query','');openSearch()}break}

  case'toSelected':if(selected){cancelRoute();stopFollowing();setHasStart(false);p.setField('startQuery','');setStops([selected]);p.setField('endQuery',selected.name);setTab('directions');setSelected(null);setRoutes([]);setExpanded(true)}break;

  case'start':if(selected){cancelRoute();stopFollowing();if(!hasStart&&stops.length>=8)throw Error('Use up to eight stops per route');setHasStart(true);setStops(s=>[selected,...(hasStart?s.slice(1):s)]);p.setField('startQuery',selected.name);setRoutes([]);setSelected(null);setTab('directions');setExpanded(true)}break;

  case'stop':if(selected){cancelRoute();stopFollowing();if(stops.length>=8)throw Error('Use up to eight stops per route');setStops(s=>[...s,selected]);p.setField('endQuery',selected.name);setRoutes([]);setSelected(null);setTab('directions');setExpanded(true)}break;

  case'findStart':case'findEnd':{const field=a==='findStart'?'startQuery':'endQuery',target=a==='findStart'?'start':'end';if(!f[field]||f[field].trim().length<2)throw Error('Enter a place or address');const stamp=++searchGeneration.current;p.setBusy('Finding route stop…');const r=await p.request('/api/search?q='+encodeURIComponent(f[field])+'&'+coords());if(!alive.current||stamp!==searchGeneration.current)break;setSearchTarget(target);setRouteResults(mapFeatures(r.features));if(!r.features?.length)p.setNotice('No places found');break}

  case'routePlace':{cancelRoute();stopFollowing();const x=routeResults.find(x=>x.key===id);if(x){if(searchTarget==='start'){setHasStart(true);setStops(s=>[x,...(hasStart?s.slice(1):s)]);p.setField('startQuery',x.name)}else{setStops(s=>hasStart&&s.length?[...s.slice(0,1),x]:[x]);p.setField('endQuery',x.name)}setRoutes([]);setRouteResults([]);setSelected(null);mapRef.current?.easeTo({center:[x.lon,x.lat],zoom:14})}break}

  case'mode':{const mode=['car','bike','foot'].includes(id)?id:'car';cancelRoute();stopFollowing();p.setField('mode',mode);p.setField('vehicle',defaultVehicle(mode));setRoutes([]);if(stops.length>=2)await calculateRoute(stops,{mode,vehicle:defaultVehicle(mode)});break}case'removeStop':{cancelRoute();stopFollowing();const next=stops.filter((_,i)=>String(i)!==id),start=hasStart&&(id!=='0'||next.length>1);setStops(next);setHasStart(start);p.setField('startQuery',start?next[0]?.name||'':'');p.setField('endQuery',next.length>(start?1:0)?next.at(-1)?.name||'':'');setRoutes([]);break}case'clearStops':cancelRoute();stopFollowing();setHasStart(false);setStops([]);setRoutes([]);p.setField('startQuery','');p.setField('endQuery','');break;case'swapStops':{cancelRoute();stopFollowing();const next=[...stops].reverse();setStops(next);setHasStart(next.length>1);p.setField('startQuery',next.length>1?next[0].name:'');p.setField('endQuery',next.at(-1)?.name||'');setRoutes([]);break}

  case'route':if(stops.length<2||!hasStart)throw Error('Choose a starting point and a destination');await calculateRoute(stops);break;

  case'cancelRoute':cancelRoute();break;

  case'routeOptions':setOptionsOpen(v=>!v);setExpanded(true);break;

  case'addStops':setMainTab('explore');p.setNotice('Choose another place on the map, then tap Add stop.');break;

  case'useLocation':{cancelRoute();const stamp=routeEpoch.current;p.setBusy('Finding your location...');const start=await currentLocation();if(alive.current&&stamp===routeEpoch.current){setStops(s=>[start,...(hasStart?s.slice(1):s)]);setHasStart(true);p.setField('startQuery',start.name);setRoutes([])}break}

  case'reroute':if(navProgress?.location){const list=[navProgress.location,...stops.slice(navStopIndex.current)];setStops(list);if(await calculateRoute(list))locate(true)}break;

  case'chooseRoute':stopFollowing();setRouteIndex(Number(id));routeCurrent.current=routes[Number(id)];break;

  case'shareRoute':{if(stops.length<2)break;const u=new URL(location.href);u.hash=new URLSearchParams({route:JSON.stringify({mode:f.mode||'car',stops:stops.map(x=>({name:x.name.slice(0,100),lat:x.lat,lon:x.lon}))})}).toString();if(navigator.share){try{await navigator.share({title:'Directions to '+stops.at(-1)?.name,url:u.href})}catch(e:any){if(e.name!=='AbortError')throw e}}else{await navigator.clipboard.writeText(u.href);p.setNotice('Route link copied. Only stop names and coordinates are shared.')}break}

  case'save':if(selected)openSave(selected);break;

  case'saveResult':{const x=results.find(x=>x.key===id);if(x){if(shortcutIntent)choose(x);else openSave(x)}break}

  case'saveSuggestion':{const x=suggestions.find(x=>x.key===id);if(x){if(shortcutIntent)choose(x);else openSave(x)}break}

  case'saveRecent':{const x=recents[Number(id)];if(x&&Number.isFinite(x.lat)&&Number.isFinite(x.lon)){if(shortcutIntent)choose(x);else openSave(x)}break}

  case'editSaved':{const x=places.find(x=>x.id===id);if(x)openSave(x);break}

  case'confirmSave':if(saveTarget){

   const choice=String(f.saveList||'Favourites'),typed=String(f.newList||'').trim();

   let list=choice==='__new_list__'?typed:choice;

   if(!list||new TextEncoder().encode(list).length>100)throw Error('Enter a short list name before saving.');

   if(list==='__new_list__')throw Error('Choose a different list name.');

   list=listNames.find(x=>x.toLocaleLowerCase()===list.toLocaleLowerCase())||list;

   p.setBusy('Saving place...');const current=saveTarget;

   const x=await p.request('/api/places',{method:'POST',body:JSON.stringify({id:current.id||'',name:current.name,address:current.address,lat:current.lat,lon:current.lon,list,notes:f.saveNotes||''})});

   if(alive.current){if(selected&&samePlace(selected,current)){setSelected({...selected,...x});p.setField('list',x.list);p.setField('notes',x.notes)}const updated=await loadPlaces();setSelected((prev:Place|null)=>{if(!prev)return prev;const fresh=updated.find((v:Place)=>samePlace(v,prev));return fresh?{...prev,...fresh}:prev});close();p.setNotice('Saved to '+x.list)}

  }break;

  case'deletePlace':await p.request('/api/places/'+encodeURIComponent(id),{method:'DELETE'});await loadPlaces();if(selected?.id===id)setSelected(null);break;

  case'share':if(selected){const u=new URL(location.href);u.hash=new URLSearchParams({lat:String(selected.lat),lon:String(selected.lon),name:selected.name}).toString();if(navigator.share){try{await navigator.share({title:selected.name,url:u.href})}catch(e:any){if(e.name!=='AbortError')throw e}}else{await navigator.clipboard.writeText(u.href);p.setNotice('Place link copied. Private notes are excluded.')}}break;

  case'locate':locate();break;case'measure':close();setMeasuring(!measuring);setMeasurePoints([]);setExpanded(false);break;

  case'navigate':if(navigating)stopFollowing();else if(routes.length)locate(true);break;

  case'exportPlaces':await p.download(new Blob([JSON.stringify({type:'FeatureCollection',features:places.map(x=>({type:'Feature',properties:{name:x.name,address:x.address,list:x.list,notes:x.notes},geometry:{type:'Point',coordinates:[x.lon,x.lat]}}))},null,2)],{type:'application/geo+json'}),'my-places.geojson');break;

  case'exportRoute':{const r=routes[routeIndex];if(r)await p.download(new Blob(['<?xml version="1.0"?><gpx version="1.1" creator="My Maps" xmlns="http://www.topografix.com/GPX/1/1"><trk><name>My Maps route</name><trkseg>'+r.geometry.coordinates.map((c:number[])=>`<trkpt lat="${c[1]}" lon="${c[0]}"/>`).join('')+'</trkseg></trk></gpx>'],{type:'application/gpx+xml'}),'route.gpx');break}

 }}catch(e:any){if(alive.current&&e.name!=='AbortError'){p.setError(e.message)}}finally{if(alive.current)p.setBusy('')}};

 if(p.route==='#settings')return p.render(p.settings(),{showHeader:true},action);

 if(config&&(!config.ready||overlay==='install'))return React.createElement(MapInstall,{...p,config,centre:center(),onClose:async()=>{const next=await p.request('/api/config');if(next.installed?.id!==config.installed?.id)camera.current=null;setConfig(next);close()}});

 const r=routes[routeIndex],c=center();

 const steps=r?.legs.flatMap((l:any)=>l.steps).map((s:any)=>({icon:turnIcon(s),instruction:s.instruction||(s.maneuver.type==='depart'?'Start':s.maneuver.type==='arrive'?'Arrive':s.maneuver.type.replaceAll('_',' '))+(s.maneuver.modifier?' '+s.maneuver.modifier:'')+(s.name?' onto '+s.name:''),distance:fmtDistance(s.distance)}))||[];

 const exploring=tab==='explore'&&!selected,saving=tab==='saved'&&!selected,routing=tab==='directions',contributing=tab==='contribute'&&!selected;

 const localSuggestions=places.filter(x=>String(f.query||'').trim().length>=2&&(x.name+' '+x.address).toLowerCase().includes(String(f.query).toLowerCase())).slice(0,5);

 const decorated=results.map(x=>({...x,icon:iconFor(x),categoryLabel:categoryName(x),art:'place-art-'+iconFor(x).toLowerCase(),meta:fmtDistance(distance([c.lon,c.lat],[x.lon,x.lat]))+' away'}));

 const destinationList=String(f.saveList)==='__new_list__'?String(f.newList||'').trim():String(f.saveList||'Favourites'),oldShortcut=places.find(x=>x.list?.trim().toLowerCase()===destinationList.toLowerCase()&&saveTarget&&!samePlace(x,saveTarget));

 const mode=f.mode||'car',vehicle=f.vehicle||defaultVehicle(mode),vehicles=mode==='car'?[{value:'car',label:'Petrol / diesel car'},{value:'electric',label:'Electric car'},{value:'hybrid',label:'Hybrid car'}]:mode==='bike'?[{value:'touring',label:'Touring bicycle'},{value:'road',label:'Road bicycle'},{value:'gravel',label:'Gravel bicycle'},{value:'mountain',label:'Mountain bike'}]:[{value:'walking',label:'Walking'},{value:'hiking',label:'Hiking'}];

 const data={planning:routing&&!navigating,routePreview:routing&&!!r&&!navigating&&!expanded,routePreviewTime:r?duration(r.duration)+' · '+fmtDistance(r.distance):'',routePreviewSummary:stops.at(-1)?.name||'',needsDirections:!config?.routingAvailable||needsCoverage,routeBusy,cannotRoute:!hasStart||stops.length<2||routeBusy,optionsOpen,driving:mode==='car',vehicles,vehicleLabel:mode==='car'?'Vehicle type':mode==='bike'?'Bicycle type':'Walking profile',preferences:[{value:'fastest',label:'Fastest route'},{value:'efficient',label:'Lower energy use'}],navigationActive:navigating,navigationIcon:navProgress?.icon||'Navigation',nextInstruction:navProgress?.instruction||'Locating you...',remainingLabel:navProgress?.remaining||'Keep this page open',offRoute:!!navProgress?.offRoute,arrivalLabel:r?'Arrive around '+new Date(Date.now()+r.duration*1000).toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'})+' · estimate without live traffic':'',saveOpen:overlay==='save'&&!!saveTarget,saveTarget,saveError:overlay==='save'?p.error:'',saveBusy:overlay==='save'&&!!p.busy,saveOptions:[...listNames.map(x=>({value:x,label:x})),{value:'__new_list__',label:'Create a new list…'}],newListSelected:f.saveList==='__new_list__',saveConfirmLabel:destinationList==='Home'?'Save as Home':destinationList==='Work'?'Save as Work':saveTarget?.id?'Save changes':'Save place',shortcutReplacement:oldShortcut&&['home','work'].includes(destinationList.toLowerCase())?'Your current '+destinationList+' ('+oldShortcut.name+') will move to Favourites.':'',shortcutHint:shortcutIntent?'Choose a location to save as '+shortcutIntent+'.':'',searchHeading:shortcutIntent?'Set '+shortcutIntent:'',exploringNearby:exploring&&!!config?.nearbyAvailable,overviewExploring:exploring&&!config?.nearbyAvailable,searchOpen,searchHasQuery:!!String(f.query||'').trim(),searchPending:suggestionPending,searchStatus,suggestions:suggestions.map(x=>({...x,saveAria:'Save '+x.name+' to a list',saveIcon:places.some(v=>samePlace(v,x))?'BookmarkCheck':'Bookmark',saveLabel:places.some(v=>samePlace(v,x))?'Saved':'Save'})),localSuggestions,hasSuggestions:!!suggestions.length,showRecent:!String(f.query||'').trim(),recents:recents.map((x,i)=>({...x,key:String(i),title:x.name||x.query,saveable:Number.isFinite(x.lat)&&Number.isFinite(x.lon),saveAria:'Save '+x.name+' to a list'})),hasRecents:!!recents.length,homeLabel:places.find(x=>x.list?.trim().toLowerCase()==='home')?.name||'Set home',workLabel:places.find(x=>x.list?.trim().toLowerCase()==='work')?.name||'Set work',searchText:(selected?(selected.address.startsWith(selected.name)?selected.address:selected.name):'')||(category?categories.find(x=>x.key===category)?.title:'')||(attempted?resultTitle:'')||'Search here',showSearchReset:!!selected||attempted,showTabs:!selected&&!attempted&&tab!=='directions',showPlace:!!selected,placeSaved:!!selected?.id,saveLabel:selected?.id?'Saved':'Save',exploring,saving,routing,contributing,initial:(p.user.name||p.user.email||'M').slice(0,1).toUpperCase(),userName:p.user.name||'Your account',userEmail:p.user.email,workspace:p.cfg.workspace,logo:new URL('./logo.svg',location.href).pathname,versionText:'My Maps · '+p.cfg.version,title:routing?(navigating?'Navigation':mode==='car'?'Drive':mode==='bike'?'Cycle':'Walk'):selected?.name|| (saving?'You':contributing?'Contribute':category?categories.find(x=>x.key===category)?.title||'Nearby places':resultTitle||'Local vibe'),profileOpen:overlay==='profile',coordinatesOpen:overlay==='coordinates',layersOpen:overlay==='layers',layerMessage,

  attempted,emptyMessage,sheetExpanded:expanded,sheetHint:routing?(expanded||navigating?'':r?stops.length+' stops · Drag up for directions':'Choose a start and destination'):expanded?(attempted?results.length+' results':''):(selected?'':attempted?results.length+' places':(config?.nearbyAvailable?'Drag up to explore':'Offline map · Drag up for your places')),sheetVariant:'sheet-'+(expanded?'expanded':'peek'),expandLabel:expanded?'Collapse place sheet':'Expand place sheet',seeMoreLabel:expanded?'Show more map':'See more',closable:!!selected,hasResults:!!results.length,resultCaption:category?'Near this map · '+results.length+' places':results.length+' search results',results:decorated.map(x=>({...x,saveAria:'Save '+x.name+' to a list',saveIcon:places.some(v=>samePlace(v,x))?'BookmarkCheck':'Bookmark',saveLabel:places.some(v=>samePlace(v,x))?'Saved':'Save'})),categories:(config?.nearbyAvailable?(selected?[categories[0],categories[3],categories[2],...categories.filter(x=>!['restaurants','attractions','hotels'].includes(x.key))]:categories):[]).map(x=>({...x,variant:x.key===category?'chip-active':'chip'})),exploreCategories:config?.nearbyAvailable?categories:[],nearbyAvailable:!!config?.nearbyAvailable,routingAvailable:!!config?.routingAvailable,places,noPlaces:!places.length,selected:selected?{...selected,website:safeURL(selected.website),callURL:'tel:'+selected.phone?.replace(/[^+\d]/g,''),categoryLabel:categoryName(selected),coordinates:selected.lat.toFixed(5)+', '+selected.lon.toFixed(5)}:null,

  stops:stops.map((s,i)=>({key:String(i),label:(!hasStart&&stops.length===1?'Destination':String(i+1))+'. '+s.name})),notEnoughStops:stops.length<2,routeResults,modes:[{key:'car',label:'Drive',icon:'Car'},{key:'foot',label:'Walk',icon:'Footprints'},{key:'bike',label:'Cycle',icon:'Bike'}].map(x=>({...x,active:mode===x.key,variant:mode===x.key?'mode-active':'mode'})),hasRoutes:!!routes.length,routes:routes.map((r,i)=>({key:String(i),active:i===routeIndex,variant:i===routeIndex?'route-active':'route',time:duration(r.duration),summary:fmtDistance(r.distance)+(i===0?' · Recommended':' · Alternative')})),steps,navigationLabel:navigating?'Exit':'Start',measuring,measurement:'Distance: '+fmtDistance(measurePoints.reduce((n,x,i)=>i?n+distance(measurePoints[i-1],x):0,0)),

  exploreVariant:tab==='explore'?'tab-active':'tab',savedVariant:tab==='saved'?'tab-active':'tab',contributeVariant:tab==='contribute'?'tab-active':'tab',editURL:'https://www.openstreetmap.org/edit#map=17/'+c.lat+'/'+c.lon,

  mapTypes:[],mapTypeCaption:'Installed local map',detailLayers:[]};

 void moved;

 return p.render(React.createElement(Maps,{data,onAction:action}),{showHeader:false,context:{fields:{...f,vehicle,preference:f.preference||'fastest',alternatives:f.alternatives!==false},setField:(name:string,value:any)=>{p.setField(name,value);if(['vehicle','preference','avoidTolls','avoidMotorways','avoidUnpaved','alternatives','startQuery','endQuery'].includes(name)){cancelRoute();stopFollowing();setRoutes([]);routeCurrent.current=null;if(name==='startQuery'){setHasStart(false);setStops(s=>hasStart?s.slice(1):s)}if(name==='endQuery'){setStops(s=>hasStart?s.slice(0,1):[])}}if(overlay==='save')p.setError('')},dismiss:()=>{if(overlay==='save'&&p.busy)return;close()},sheetDrag:(delta:number)=>setExpanded(delta<0)}},action,{'directions-download':React.createElement(DirectionsInstall,{request:p.request,owner:config?.canInstall,status:config?.directions,centre:center(),onReady:(value:any)=>{setConfig((c:any)=>({...c,directions:value,routingAvailable:true}));setNeedsCoverage(false)}}),map:config?React.createElement(MapView,{mapRef,camera,config,base:new URL('.',location.href).pathname.replace(/\/$/,''),onPoint:pick,onChoose:choose,onMove:()=>{const m=mapRef.current;if(m)camera.current={center:m.getCenter().toArray(),zoom:m.getZoom(),bearing:m.getBearing(),pitch:m.getPitch()};setMoved(v=>v+1)},onError:p.setError,onReady:()=>{if(markerMap.current!==mapRef.current){locationMarker.current?.remove();locationMarker.current=null;markerMap.current=mapRef.current}setGeneration(v=>v+1)},generation,viewportGeneration:moved,places:tab==='saved'?places:results,stops:routing?stops:[],selected,route:r,measuring,measurePoints}):null});

}
