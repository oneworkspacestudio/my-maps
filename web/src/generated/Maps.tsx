// Generated from screens/Maps.aui; edit the AUI source.
import { Alert, Button, Card, Divider, Field, Grid, Heading, Icon, Image, Input, Link, Page, Row, Section, Spacer, Stack, Text } from "@/components/ui"

export function Maps({ data, onAction }: { data: any; onAction: (name: string) => void }) {
  return (
    <Page>
      <Stack variant="map-stage">
        <Section id="map" />
        <Stack variant="floating-search">
          <Row variant="search-bar">
            <Icon name="MapPin" size={28} />
            <Input name="query" placeholder="Search here" onClick={() => onAction("search")} />
            <Button variant="search-icon" onClick={() => onAction("voice")} aria-label="Search by voice"><Icon name="Mic" size={24} /></Button>
            <Button variant="search-icon" onClick={() => onAction("coordinates")} aria-label="Search by coordinates"><Icon name="ScanLine" size={24} /></Button>
            <Button variant="avatar" onClick={() => onAction("profile")} aria-label="Account and settings"><Text>{data.initial}</Text></Button>
          </Row>
          <Row variant="category-chips">
            <Button variant="chip" onClick={() => onAction("home")}>
              <Icon name="House" />
              <Text>Home</Text>
            </Button>
            {data.categories.map((item, i) => (
              <Button key={i} variant={item.variant} onClick={() => onAction("category")} item={item.key}>
                <Icon name={item.icon} />
                <Text>{item.title}</Text>
              </Button>
            ))}
          </Row>
        </Stack>
        <Button variant="layers-control" onClick={() => onAction("layers")} aria-label="Map layers"><Icon name="Layers2" size={26} /></Button>
        <Stack variant="map-controls">
          <Button variant="map-control" onClick={() => onAction("compass")} aria-label="Reset map bearing"><Icon name="Compass" size={25} /></Button>
          <Button variant="map-control" onClick={() => onAction("locate")} aria-label="My location"><Icon name="LocateFixed" size={25} /></Button>
          <Button variant="directions-control" onClick={() => onAction("directions")} aria-label="Directions"><Icon name="Navigation" size={26} /></Button>
        </Stack>
        <Text variant="map-brand">My Maps</Text>
        {data.measuring && (
          <Row variant="measure-status">
            <Text>{data.measurement}</Text>
            <Button variant="text" onClick={() => onAction("measure")}>Done</Button>
          </Row>
        )}
        <Stack variant={data.sheetVariant}>
          <Button variant="sheet-handle" onClick={() => onAction("expand")} aria-label={data.expandLabel}><Section variant="handle-line" /></Button>
          <Row variant="sheet-heading">
            <Heading level={1}>{data.title}</Heading>
            <Spacer />
            {data.exploring && (
              <Button variant="weather" onClick={() => onAction("weather")} aria-label="Local weather">
                <Icon name="CloudSun" />
                <Text>{data.temperature}</Text>
              </Button>
            )}
            {data.closable && (
              <Button variant="search-icon" onClick={() => onAction("deselect")} aria-label="Close place details"><Icon name="X" /></Button>
            )}
          </Row>
          <Stack variant="sheet-content">
            {data.exploring && (
              <>
                {data.hasResults ? (
                  <>
                    <Row variant="results-heading">
                      <Text variant="caption">{data.resultCaption}</Text>
                      <Spacer />
                      <Button variant="text" onClick={() => onAction("refresh")}>Search this area</Button>
                    </Row>
                    <Row variant="discovery-cards">
                      {data.results.map((item, i) => (
                        <Button key={i} variant="discovery-card" onClick={() => onAction("place")} item={item.key}>
                          <Stack variant={item.art}>
                            <Icon name={item.icon} size={52} />
                            <Text variant="card-category">{item.categoryLabel}</Text>
                          </Stack>
                          <Stack variant="card-copy">
                            <Heading level={2}>{item.name}</Heading>
                            <Text variant="caption">{item.address}</Text>
                            <Text variant="place-meta">{item.meta}</Text>
                          </Stack>
                        </Button>
                      ))}
                    </Row>
                    <Button variant="see-more" onClick={() => onAction("expand")}>{data.seeMoreLabel}</Button>
                  </>
                ) : (
                  <Card variant="discover-intro">
                    <Icon name="Utensils" size={32} />
                    <Heading level={2}>Find your local favourites</Heading>
                    <Text>Discover somewhere to eat, shop or explore around this map.</Text>
                    <Button variant="pill" onClick={() => onAction("restaurants")}>Find restaurants</Button>
                  </Card>
                )}
                <Heading level={2}>Explore nearby</Heading>
                <Grid variant="explore-grid" columns={2}>
                  {data.exploreCategories.map((item, i) => (
                    <Button key={i} variant="explore-tile" onClick={() => onAction("category")} item={item.key}>
                      <Icon name={item.icon} />
                      <Text>{item.title}</Text>
                    </Button>
                  ))}
                </Grid>
                <Text variant="caption">Place details from OpenStreetMap. Coverage varies by area.</Text>
              </>
            )}
            {data.saving && (
              <>
                <Text>Your places, all together. Save a place on the map to build your own lists.</Text>
                <Row variant="saved-shortcuts">
                  <Button variant="shortcut" onClick={() => onAction("home")}>
                    <Icon name="House" />
                    <Text>Home</Text>
                  </Button>
                  <Button variant="shortcut" onClick={() => onAction("work")}>
                    <Icon name="BriefcaseBusiness" />
                    <Text>Work</Text>
                  </Button>
                  <Button variant="shortcut" onClick={() => onAction("exportPlaces")}>
                    <Icon name="Download" />
                    <Text>Export</Text>
                  </Button>
                </Row>
                {data.noPlaces && (
                  <Card variant="empty-places">
                    <Icon name="Bookmark" size={36} />
                    <Heading level={2}>Keep the places you love</Heading>
                    <Text>Choose a place and tap Save to add it to a list.</Text>
                  </Card>
                )}
                {data.places.map((item, i) => (
                  <Row key={i} variant="saved-place">
                    <Button variant="result" onClick={() => onAction("savedPlace")} item={item.id}>
                      <Icon name="MapPin" />
                      <Stack>
                        <Text>{item.name}</Text>
                        <Text variant="caption">{item.list}</Text>
                      </Stack>
                    </Button>
                    <Button variant="search-icon" onClick={() => onAction("deletePlace")} item={item.id} aria-label="Delete saved place"><Icon name="Trash2" /></Button>
                  </Row>
                ))}
              </>
            )}
            {data.contributing && (
              <>
                <Card variant="contribute-intro">
                  <Icon name="MapPinned" size={42} />
                  <Heading level={2}>Help make the map better</Heading>
                  <Text>Add a missing place or correct details with the OpenStreetMap community.</Text>
                  <Link variant="button" href={data.editURL} external={true}>Edit this area</Link>
                  <Link href="https://www.openstreetmap.org/fixthemap" external={true}>Report a map problem</Link>
                </Card>
                <Heading level={2}>Your place notes</Heading>
                <Text>Choose a place on the map, save it, and add your own private notes.</Text>
                <Button variant="pill" onClick={() => onAction("explore")}>Choose a place</Button>
              </>
            )}
            {data.routing && (
              <>
                <Row variant="travel-modes">
                  {data.modes.map((item, i) => (
                    <Button key={i} variant={item.variant} onClick={() => onAction("mode")} item={item.key}>
                      <Icon name={item.icon} />
                      <Text>{item.label}</Text>
                    </Button>
                  ))}
                </Row>
                <Row>
                  <Icon name="CircleDot" />
                  <Input name="startQuery" placeholder="Choose starting point" onClick={() => onAction("findStart")} />
                  <Button variant="search-icon" onClick={() => onAction("findStart")} aria-label="Find starting point"><Icon name="Search" /></Button>
                </Row>
                <Row>
                  <Icon name="MapPin" />
                  <Input name="endQuery" placeholder="Choose destination" onClick={() => onAction("findEnd")} />
                  <Button variant="search-icon" onClick={() => onAction("findEnd")} aria-label="Find destination"><Icon name="Search" /></Button>
                </Row>
                {data.routeResults.map((item, i) => (
                  <Button key={i} variant="result" onClick={() => onAction("routePlace")} item={item.key}>
                    <Text>{item.name}</Text>
                    <Text variant="caption">{item.address}</Text>
                  </Button>
                ))}
                {data.stops.map((item, i) => (
                  <Row key={i} variant="stop">
                    <Text>{item.label}</Text>
                    <Spacer />
                    <Button variant="search-icon" onClick={() => onAction("removeStop")} item={item.key} aria-label="Remove stop"><Icon name="X" /></Button>
                  </Row>
                ))}
                <Row>
                  <Button variant="pill" onClick={() => onAction("route")} disabled={data.notEnoughStops}>Find routes</Button>
                  <Button variant="text" onClick={() => onAction("swapStops")}>Swap</Button>
                  <Button variant="text" onClick={() => onAction("clearStops")}>Clear</Button>
                </Row>
                <Text variant="caption">Or choose places on the map, then tap Set start or Add stop.</Text>
                {data.hasRoutes && (
                  <>
                    {data.routes.map((item, i) => (
                      <Button key={i} variant={item.variant} onClick={() => onAction("chooseRoute")} item={item.key}>
                        <Icon name="Route" />
                        <Text>{item.summary}</Text>
                      </Button>
                    ))}
                    <Row>
                      <Button variant="pill" onClick={() => onAction("navigate")}>{data.navigationLabel}</Button>
                      <Button variant="secondary" onClick={() => onAction("exportRoute")}>Export GPX</Button>
                    </Row>
                    <Text variant="caption">Route estimates exclude live traffic and road closures. Follow-location works while this page is open.</Text>
                    {data.steps.map((item, i) => (
                      <Row key={i} variant="step">
                        <Icon name="CornerUpRight" />
                        <Text>{item.instruction}</Text>
                        <Spacer />
                        <Text variant="caption">{item.distance}</Text>
                      </Row>
                    ))}
                  </>
                )}
              </>
            )}
            {data.selected && (
              <Card variant="place-detail">
                <Text variant="place-category">{data.selected.categoryLabel}</Text>
                <Text>{data.selected.address}</Text>
                <Row variant="place-actions">
                  <Button variant="pill" onClick={() => onAction("toSelected")}>
                    <Icon name="Navigation" />
                    <Text>Directions</Text>
                  </Button>
                  <Button variant="outline-pill" onClick={() => onAction("save")}>
                    <Icon name="Bookmark" />
                    <Text>Save</Text>
                  </Button>
                  <Button variant="outline-pill" onClick={() => onAction("share")}>
                    <Icon name="Share2" />
                    <Text>Share</Text>
                  </Button>
                </Row>
                <Divider />
                <Row>
                  <Icon name="MapPin" />
                  <Text>{data.selected.coordinates}</Text>
                </Row>
                {data.selected.hours && (
                  <Row>
                    <Icon name="Clock" />
                    <Text>{data.selected.hours}</Text>
                  </Row>
                )}
                {data.selected.website && (
                  <Link href={data.selected.website} external={true}>Visit website</Link>
                )}
                {data.selected.phone && (
                  <Link href={data.selected.callURL}>{data.selected.phone}</Link>
                )}
                {data.selected.osm && (
                  <Link href={data.selected.osm} external={true}>Place details on OpenStreetMap</Link>
                )}
                <Row>
                  <Button variant="text" onClick={() => onAction("start")}>Set start</Button>
                  <Button variant="text" onClick={() => onAction("stop")}>Add stop</Button>
                  <Button variant="text" onClick={() => onAction("measure")}>Measure distance</Button>
                </Row>
                <Field label="List"><Input name="list" placeholder="Favourites" /></Field>
                <Field label="Your private notes"><Input name="notes" multiline={true} /></Field>
                <Button variant="pill" onClick={() => onAction("save")}>Save place and notes</Button>
              </Card>
            )}
          </Stack>
        </Stack>
        <Row variant="bottom-tabs">
          <Button variant={data.exploreVariant} onClick={() => onAction("explore")}>
            <Icon name="MapPin" size={25} />
            <Text>Explore</Text>
          </Button>
          <Button variant={data.savedVariant} onClick={() => onAction("saved")}>
            <Icon name="Bookmark" size={25} />
            <Text>You</Text>
          </Button>
          <Button variant={data.contributeVariant} onClick={() => onAction("contribute")}>
            <Icon name="CirclePlus" size={25} />
            <Text>Contribute</Text>
          </Button>
        </Row>
      </Stack>
      {data.layersOpen && (
        <Card variant="dialog" aria-label="Map layers">
          <Row>
            <Heading level={2}>Map type</Heading>
            <Spacer />
            <Button variant="search-icon" onClick={() => onAction("closeOverlay")} cancel={true} aria-label="Close map layers"><Icon name="X" size={26} /></Button>
          </Row>
          <Grid variant="map-types" columns={3}>
            {data.mapTypes.map((item, i) => (
              <Button key={i} variant={item.variant} onClick={() => onAction("mapType")} item={item.key} aria-pressed={item.active}>
                <Section variant={item.art}><Icon name={item.icon} size={42} /></Section>
                <Text>{item.title}</Text>
              </Button>
            ))}
          </Grid>
          <Text variant="caption">{data.mapTypeCaption}</Text>
          <Divider />
          <Heading level={2}>Map details</Heading>
          <Grid variant="map-details" columns={4}>
            {data.detailLayers.map((item, i) => (
              <Button key={i} variant={item.variant} onClick={() => onAction("detailLayer")} item={item.key} aria-pressed={item.active}>
                <Section variant={item.art}><Icon name={item.icon} size={38} /></Section>
                <Text>{item.title}</Text>
              </Button>
            ))}
          </Grid>
          {data.layerMessage && (
            <Alert>{data.layerMessage}</Alert>
          )}
          {data.layerURL && (
            <Link variant="button-secondary" href={data.layerURL} external={true}>{data.layerLink}</Link>
          )}
        </Card>
      )}
      {data.profileOpen && (
        <Card variant="dialog" aria-label="Account and settings">
          <Row>
            <Image src={data.logo} alt="My Maps" width={32} height={32} />
            <Heading level={2}>My Maps</Heading>
            <Spacer />
            <Button variant="search-icon" onClick={() => onAction("closeOverlay")} aria-label="Close account menu"><Icon name="X" /></Button>
          </Row>
          <Row variant="profile-identity">
            <Icon name="CircleUserRound" size={40} />
            <Stack>
              <Heading level={2}>{data.userName}</Heading>
              <Text variant="caption">{data.userEmail}</Text>
            </Stack>
          </Row>
          <Button variant="menu-item" onClick={() => onAction("saved")}>
            <Icon name="Bookmark" />
            <Text>Your places</Text>
          </Button>
          <Button variant="menu-item" onClick={() => onAction("measure")}>
            <Icon name="Ruler" />
            <Text>Measure distance</Text>
          </Button>
          <Button variant="menu-item" onClick={() => onAction("theme")}>
            <Icon name="SunMoon" />
            <Text>Change appearance</Text>
          </Button>
          <Link variant="menu-item" href="#settings">
            <Icon name="Settings2" />
            <Text>Settings and backups</Text>
          </Link>
          {data.workspace && (
            <Button variant="menu-item" onClick={() => onAction("apps")} data-workspace-apps-trigger={true}>
              <Icon name="Layers" />
              <Text>Workspace apps</Text>
            </Button>
          )}
          <Divider />
          <Button variant="menu-item" onClick={() => onAction("logout")}>
            <Icon name="LogOut" />
            <Text>Sign out</Text>
          </Button>
          <Text variant="caption">{data.versionText}</Text>
        </Card>
      )}
      {data.coordinatesOpen && (
        <Card variant="dialog" aria-label="Find coordinates">
          <Row>
            <Heading level={2}>Find coordinates</Heading>
            <Spacer />
            <Button variant="search-icon" onClick={() => onAction("closeOverlay")} aria-label="Close coordinates"><Icon name="X" /></Button>
          </Row>
          <Text>Enter latitude and longitude to open an exact location.</Text>
          <Field label="Coordinates"><Input name="coordinates" placeholder="50.7829, -0.6731" onClick={() => onAction("goCoordinates")} /></Field>
          <Button variant="pill" onClick={() => onAction("goCoordinates")}>Show on map</Button>
        </Card>
      )}
      {data.weatherOpen && (
        <Card variant="dialog" aria-label="Local conditions">
          <Row>
            <Heading level={2}>Local conditions</Heading>
            <Spacer />
            <Button variant="search-icon" onClick={() => onAction("closeOverlay")} aria-label="Close conditions"><Icon name="X" /></Button>
          </Row>
          <Text>{data.conditions}</Text>
          <Text variant="caption">Forecast and air-quality models from Open-Meteo and CAMS. Not a local sensor reading.</Text>
          <Link href="https://open-meteo.com/" external={true}>Weather source and attribution</Link>
        </Card>
      )}
    </Page>
  )
}
