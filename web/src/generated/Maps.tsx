// Generated from screens/Maps.aui; edit the AUI source.
import { Alert, Button, Card, Checkbox, Divider, Field, Grid, Heading, Icon, Image, Input, Link, Page, Row, Section, Select, Spacer, Spinner, Stack, Text } from "@/components/ui"

export function Maps({ data, onAction }: { data: any; onAction: (name: string) => void }) {
  return (
    <Page>
      <Stack variant="map-stage">
        <Section id="map" />
        <Stack variant="floating-search">
          <Row variant="search-bar">
            <Icon name="MapPin" size={28} />
            <Button variant="search-entry" onClick={() => onAction("openSearch")} aria-label="Search places"><Text>{data.searchText}</Text></Button>
            {data.showSearchReset ? (
              <Button variant="search-icon" onClick={() => onAction("clearSearch")} aria-label="Clear search"><Icon name="X" size={24} /></Button>
            ) : (
              <Button variant="avatar" onClick={() => onAction("profile")} aria-label="Account and settings"><Text>{data.initial}</Text></Button>
            )}
          </Row>
          <Row variant="category-chips">
            {data.showTabs && (
              <Button variant="chip" onClick={() => onAction("home")}>
                <Icon name="House" />
                <Text>Home</Text>
              </Button>
            )}
            {data.categories.map((item, i) => (
              <Button key={i} variant={item.variant} onClick={() => onAction("category")} item={item.key}>
                <Icon name={item.icon} />
                <Text>{item.title}</Text>
              </Button>
            ))}
          </Row>
        </Stack>
        <Button variant="layers-control" onClick={() => onAction("layers")} aria-label="Map downloads"><Icon name="Download" size={26} /></Button>
        <Stack variant="map-controls">
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
          <Button variant="sheet-handle" onClick={() => onAction("expand")} aria-label={data.expandLabel} aria-expanded={data.sheetExpanded}><Section variant="handle-line" /></Button>
          <Row variant="sheet-heading">
            {data.routing && (
              <Button variant="search-icon" onClick={() => onAction("clearSearch")} aria-label="Back to map"><Icon name="ArrowLeft" /></Button>
            )}
            <Stack variant="sheet-title">
              <Heading level={1}>{data.title}</Heading>
              <Text variant="sheet-hint">{data.sheetHint}</Text>
            </Stack>
            {data.closable && (
              <>
                <Button variant="round-action" onClick={() => onAction("save")} aria-label="Save place"><Icon name="Bookmark" /></Button>
                <Button variant="round-action" onClick={() => onAction("share")} aria-label="Share place"><Icon name="Share2" /></Button>
                <Button variant="search-icon" onClick={() => onAction("deselect")} aria-label="Close place details"><Icon name="X" /></Button>
              </>
            )}
          </Row>
          {data.selected && (
            <>
              <Stack variant="place-summary">
                <Text variant="place-address">{data.selected.address}</Text>
                <Text variant="caption">{data.selected.categoryLabel}</Text>
              </Stack>
              <Row variant="quick-actions">
                <Button variant="pill" onClick={() => onAction("toSelected")}>
                  <Icon name="CornerUpRight" />
                  <Text>Directions</Text>
                </Button>
                <Button variant="soft-pill" onClick={() => onAction("quickStart")}>
                  <Icon name="Navigation" />
                  <Text>Start</Text>
                </Button>
                <Button variant="soft-pill" onClick={() => onAction("share")}>
                  <Icon name="Share2" />
                  <Text>Share</Text>
                </Button>
                <Button variant="soft-pill" onClick={() => onAction("save")}>
                  <Icon name="Bookmark" />
                  <Text>{data.saveLabel}</Text>
                </Button>
              </Row>
            </>
          )}
          {data.routePreview && (
            <Stack variant="route-preview">
              <Text variant="route-time">{data.routePreviewTime}</Text>
              <Text variant="caption">{data.routePreviewSummary}</Text>
              <Row variant="quick-actions">
                <Button variant="pill" onClick={() => onAction("navigate")}>
                  <Icon name="Navigation" />
                  <Text>{data.navigationLabel}</Text>
                </Button>
                <Button variant="soft-pill" onClick={() => onAction("addStops")}>
                  <Icon name="CirclePlus" />
                  <Text>Add stops</Text>
                </Button>
                <Button variant="soft-pill" onClick={() => onAction("shareRoute")}>
                  <Icon name="Share2" />
                  <Text>Share</Text>
                </Button>
              </Row>
            </Stack>
          )}
          {data.navigationActive && (
            <Card variant="navigation-cue">
              <Icon name={data.navigationIcon} size={28} />
              <Heading level={2}>{data.nextInstruction}</Heading>
              <Text>{data.remainingLabel}</Text>
              {data.offRoute && (
                <Button variant="pill" onClick={() => onAction("reroute")}>Update route from here</Button>
              )}
              <Button variant="soft-pill" onClick={() => onAction("navigate")}>Exit navigation</Button>
            </Card>
          )}
          <Stack variant="sheet-content">
            {data.exploringNearby && (
              <>
                {data.hasResults ? (
                  <>
                    <Row variant="results-heading">
                      <Text variant="caption">{data.resultCaption}</Text>
                      <Spacer />
                      <Button variant="text" onClick={() => onAction("refresh")}>Search this area</Button>
                    </Row>
                    <Stack variant="place-results">
                      {data.results.map((item, i) => (
                        <Row key={i} variant="place-result-row">
                          <Button variant="place-result" onClick={() => onAction("place")} item={item.key}>
                            <Icon name={item.icon} size={23} />
                            <Stack variant="result-copy">
                              <Heading level={2}>{item.name}</Heading>
                              <Text variant="caption">{item.address}</Text>
                              <Text variant="place-meta">{item.meta}</Text>
                            </Stack>
                          </Button>
                          <Button variant="list-save" onClick={() => onAction("saveResult")} item={item.key} aria-label={item.saveAria}>
                            <Icon name={item.saveIcon} size={23} />
                            <Text>{item.saveLabel}</Text>
                          </Button>
                        </Row>
                      ))}
                    </Stack>
                    <Button variant="see-more" onClick={() => onAction("expand")}>{data.seeMoreLabel}</Button>
                  </>
                ) : (
                  <Stack variant="explore-empty">
                    {data.attempted ? (
                      <>
                        <Text>{data.emptyMessage}</Text>
                        <Button variant="see-more" onClick={() => onAction("refresh")}>Try again</Button>
                      </>
                    ) : (
                      <Card variant="discover-intro">
                        <Icon name="Utensils" size={32} />
                        <Heading level={2}>Find your local favourites</Heading>
                        <Text>Discover somewhere to eat, shop or explore around this map.</Text>
                        <Button variant="pill" onClick={() => onAction("restaurants")}>Find restaurants</Button>
                      </Card>
                    )}
                  </Stack>
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
            {data.overviewExploring && (
              <Card>
                <Heading level={2}>Your local overview</Heading>
                <Text>Search towns and mapped names, save places, or install street detail for nearby places.</Text>
                <Button variant="soft-pill" onClick={() => onAction("mapInstall")}>Map downloads</Button>
              </Card>
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
                    <Button variant="search-icon" onClick={() => onAction("editSaved")} item={item.id} aria-label="Change saved place list"><Icon name="Bookmark" /></Button>
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
            {data.planning && (
              <>
                <Row variant="travel-modes">
                  {data.modes.map((item, i) => (
                    <Button key={i} variant={item.variant} onClick={() => onAction("mode")} item={item.key} aria-pressed={item.active}>
                      <Icon name={item.icon} />
                      <Text>{item.label}</Text>
                    </Button>
                  ))}
                </Row>
                {data.needsDirections && (
                  <Section id="directions-download" />
                )}
                <Stack variant="route-inputs">
                  <Row>
                    <Icon name="CircleDot" />
                    <Input name="startQuery" placeholder="Choose starting point" onClick={() => onAction("findStart")} />
                    <Button variant="search-icon" onClick={() => onAction("useLocation")} aria-label="Use my location as start"><Icon name="LocateFixed" /></Button>
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
                </Stack>
                <Row variant="route-options-heading">
                  <Button variant="option-chip" onClick={() => onAction("routeOptions")} aria-expanded={data.optionsOpen}>
                    <Icon name="Settings2" />
                    <Text>Options & vehicle</Text>
                  </Button>
                  <Button variant="text" onClick={() => onAction("swapStops")}>Swap</Button>
                  <Button variant="text" onClick={() => onAction("clearStops")}>Clear</Button>
                </Row>
                {data.optionsOpen && (
                  <Card variant="route-options">
                    <Field label={data.vehicleLabel}><Select name="vehicle" options={data.vehicles} /></Field>
                    {data.driving && (
                      <>
                        <Field label="Route preference"><Select name="preference" options={data.preferences} /></Field>
                        <Checkbox name="avoidTolls" label="Avoid tolls" />
                        <Checkbox name="avoidMotorways" label="Avoid motorways" />
                        <Checkbox name="avoidUnpaved" label="Avoid unpaved roads" />
                        <Text variant="caption">Vehicle type adjusts the local energy model. Battery range and charging stops are not calculated.</Text>
                      </>
                    )}
                    <Checkbox name="alternatives" label="Show alternative routes" />
                  </Card>
                )}
                {data.stops.map((item, i) => (
                  <Row key={i} variant="stop">
                    <Text>{item.label}</Text>
                    <Spacer />
                    <Button variant="search-icon" onClick={() => onAction("removeStop")} item={item.key} aria-label="Remove stop"><Icon name="X" /></Button>
                  </Row>
                ))}
                {data.routeBusy ? (
                  <Row variant="calculating-route">
                    <Spinner />
                    <Text>Calculating on your server...</Text>
                    <Button variant="text" onClick={() => onAction("cancelRoute")} cancel={true}>Cancel</Button>
                  </Row>
                ) : (
                  <Button variant="pill" onClick={() => onAction("route")} disabled={data.cannotRoute}>Find routes</Button>
                )}
                {data.hasRoutes ? (
                  <>
                    <Stack variant="route-choices">
                      {data.routes.map((item, i) => (
                        <Button key={i} variant={item.variant} onClick={() => onAction("chooseRoute")} item={item.key} aria-pressed={item.active}>
                          <Icon name="Route" />
                          <Stack>
                            <Text variant="route-time">{item.time}</Text>
                            <Text variant="caption">{item.summary}</Text>
                          </Stack>
                        </Button>
                      ))}
                    </Stack>
                    <Text variant="caption">{data.arrivalLabel}</Text>
                    <Row variant="route-actions">
                      <Button variant="pill" onClick={() => onAction("navigate")}>
                        <Icon name="Navigation" />
                        <Text>{data.navigationLabel}</Text>
                      </Button>
                      <Button variant="soft-pill" onClick={() => onAction("addStops")}>
                        <Icon name="CirclePlus" />
                        <Text>Add stops</Text>
                      </Button>
                      <Button variant="soft-pill" onClick={() => onAction("shareRoute")}>
                        <Icon name="Share2" />
                        <Text>Share</Text>
                      </Button>
                      <Button variant="soft-pill" onClick={() => onAction("exportRoute")}>
                        <Icon name="Download" />
                        <Text>GPX</Text>
                      </Button>
                    </Row>
                    <Text variant="caption">Estimates exclude live traffic and road closures. Navigation needs location permission and works while this page is open.</Text>
                    <Heading level={2}>Directions</Heading>
                    {data.steps.map((item, i) => (
                      <Row key={i} variant="step">
                        <Icon name={item.icon} />
                        <Text>{item.instruction}</Text>
                        <Spacer />
                        <Text variant="caption">{item.distance}</Text>
                      </Row>
                    ))}
                  </>
                ) : (
                  <Text variant="caption">Choose a start and destination, or use your current location. To add more stops, choose a place on the map and tap Add stop.</Text>
                )}
              </>
            )}
            {data.selected && (
              <Card variant="place-detail">
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
                {data.selected.list && (
                  <Text variant="caption">Saved in {data.selected.list}</Text>
                )}
                {data.selected.notes && (
                  <Text>{data.selected.notes}</Text>
                )}
                <Button variant="pill" onClick={() => onAction("save")}>Save to a list</Button>
              </Card>
            )}
          </Stack>
        </Stack>
        {data.showTabs && (
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
        )}
        {data.searchOpen && (
          <Stack variant="search-view" aria-label="Place search">
            <Row variant="search-editor">
              <Button variant="search-icon" onClick={() => onAction("closeSearch")} aria-label="Back to map"><Icon name="ArrowLeft" size={24} /></Button>
              <Input name="query" placeholder="Search here" aria-label="Search places" onClick={() => onAction("search")} autoFocus={true} enterKeyHint="search" />
              {data.searchHasQuery && (
                <Button variant="search-icon" onClick={() => onAction("clearQuery")} aria-label="Clear query"><Icon name="X" /></Button>
              )}
              <Button variant="search-icon" onClick={() => onAction("search")} aria-label="Submit search"><Icon name="Search" size={23} /></Button>
            </Row>
            {data.searchHeading && (
              <Stack variant="shortcut-intent">
                <Heading level={2}>{data.searchHeading}</Heading>
                <Text variant="caption">{data.shortcutHint}</Text>
              </Stack>
            )}
            <Row variant="search-shortcuts">
              <Button variant="search-shortcut" onClick={() => onAction("home")}>
                <Icon name="House" size={24} />
                <Stack>
                  <Text>Home</Text>
                  <Text variant="caption">{data.homeLabel}</Text>
                </Stack>
              </Button>
              <Button variant="search-shortcut" onClick={() => onAction("work")}>
                <Icon name="BriefcaseBusiness" size={24} />
                <Stack>
                  <Text>Work</Text>
                  <Text variant="caption">{data.workLabel}</Text>
                </Stack>
              </Button>
              <Button variant="search-shortcut" onClick={() => onAction("coordinates")}>
                <Icon name="LocateFixed" size={24} />
                <Text>Coordinates</Text>
              </Button>
            </Row>
            <Stack variant="search-scroll">
              {data.showRecent ? (
                <>
                  <Row variant="recent-heading">
                    <Heading level={2}>Recent</Heading>
                    <Spacer />
                    {data.hasRecents && (
                      <Button variant="text" onClick={() => onAction("clearRecents")}>Clear all</Button>
                    )}
                  </Row>
                  {data.hasRecents ? (
                    <Stack variant="recent-list">
                      {data.recents.map((item, i) => (
                        <Row key={i} variant="search-result-row">
                          <Button variant="search-result" onClick={() => onAction("recentSearch")} item={item.key}>
                            <Icon name="History" size={23} />
                            <Stack variant="result-copy">
                              <Text>{item.title}</Text>
                              <Text variant="caption">{item.address}</Text>
                            </Stack>
                          </Button>
                          {item.saveable && (
                            <Button variant="list-save" onClick={() => onAction("saveRecent")} item={item.key} aria-label={item.saveAria}>
                              <Icon name="Bookmark" size={23} />
                              <Text>Save</Text>
                            </Button>
                          )}
                        </Row>
                      ))}
                    </Stack>
                  ) : (
                    <Text variant="caption">Your searches and chosen places appear here. History stays private on this server.</Text>
                  )}
                </>
              ) : (
                <>
                  {data.localSuggestions.map((item, i) => (
                    <Button key={i} variant="search-result" onClick={() => onAction("savedPlace")} item={item.id}>
                      <Icon name="Bookmark" size={23} />
                      <Stack variant="result-copy">
                        <Text>{item.name}</Text>
                        <Text variant="caption">{item.address}</Text>
                      </Stack>
                    </Button>
                  ))}
                  {data.searchPending && (
                    <Row variant="search-status">
                      <Spinner />
                      <Text>Finding places...</Text>
                    </Row>
                  )}
                  {data.suggestions.map((item, i) => (
                    <Row key={i} variant="search-result-row">
                      <Button variant="search-result" onClick={() => onAction("suggestion")} item={item.key}>
                        <Icon name="MapPin" size={23} />
                        <Stack variant="result-copy">
                          <Text>{item.name}</Text>
                          <Text variant="caption">{item.address}</Text>
                        </Stack>
                      </Button>
                      <Button variant="list-save" onClick={() => onAction("saveSuggestion")} item={item.key} aria-label={item.saveAria}>
                        <Icon name={item.saveIcon} size={23} />
                        <Text>{item.saveLabel}</Text>
                      </Button>
                    </Row>
                  ))}
                  {data.searchStatus && (
                    <Text variant="caption">{data.searchStatus}</Text>
                  )}
                </>
              )}
            </Stack>
          </Stack>
        )}
      </Stack>
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
          <Button variant="menu-item" onClick={() => onAction("mapInstall")}>
            <Icon name="Download" />
            <Text>Map downloads</Text>
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
      {data.saveOpen && (
        <Card variant="dialog" aria-label="Save place to list">
          <Row>
            <Heading level={2}>Save to a list</Heading>
            <Spacer />
            <Button variant="search-icon" onClick={() => onAction("closeOverlay")} aria-label="Close save picker"><Icon name="X" /></Button>
          </Row>
          <Stack variant="save-place-summary">
            <Icon name="MapPin" size={24} />
            <Heading level={3}>{data.saveTarget.name}</Heading>
            <Text variant="caption">{data.saveTarget.address}</Text>
          </Stack>
          {data.saveError && (
            <Alert tone="error">{data.saveError}</Alert>
          )}
          <Field label="Choose a list"><Select name="saveList" options={data.saveOptions} /></Field>
          {data.newListSelected && (
            <Field label="New list name"><Input name="newList" placeholder="e.g. Weekend trips" maxLength={100} autoFocus={true} onClick={() => onAction("confirmSave")} /></Field>
          )}
          {data.shortcutReplacement && (
            <Text variant="save-explanation">{data.shortcutReplacement}</Text>
          )}
          <Field label="Private note (optional)"><Input name="saveNotes" multiline={true} maxLength={10000} /></Field>
          {data.saveBusy && (
            <Row>
              <Spinner />
              <Text>Saving place...</Text>
            </Row>
          )}
          <Row variant="save-picker-actions">
            <Button variant="soft-pill" onClick={() => onAction("closeOverlay")}>Cancel</Button>
            <Button variant="pill" onClick={() => onAction("confirmSave")}>{data.saveConfirmLabel}</Button>
          </Row>
        </Card>
      )}
    </Page>
  )
}
