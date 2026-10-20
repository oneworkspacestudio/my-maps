// Generated from screens/Maps.aui; edit the AUI source.
import { Button, Card, Field, Heading, Input, Link, Page, Row, Section, Select, Spacer, Stack, Text } from "@/components/ui"

export function Maps({ data, onAction }: { data: any; onAction: (name: string) => void }) {
  return (
    <Page>
      <Row variant="map-layout">
        <Stack variant="map-panel">
          <Row variant="tabs">
            <Button variant={data.exploreVariant} onClick={() => onAction("explore")}>Explore</Button>
            <Button variant={data.savedVariant} onClick={() => onAction("saved")}>Saved</Button>
            <Button variant={data.routeVariant} onClick={() => onAction("directions")}>Directions</Button>
          </Row>
          {data.exploring && (
            <>
              <Row>
                <Input name="query" placeholder="Search places and addresses" onClick={() => onAction("search")} />
                <Button onClick={() => onAction("search")} aria-label="Search map">Search</Button>
              </Row>
              <Text variant="caption">Search is sent to the community Photon service.</Text>
              {data.results.map((item, i) => (
                <Button key={i} variant="result" onClick={() => onAction("place")} item={item.key}>
                  <Stack>
                    <Text>{item.name}</Text>
                    <Text variant="caption">{item.address}</Text>
                  </Stack>
                </Button>
              ))}
            </>
          )}
          {data.saving && (
            <>
              <Row>
                <Heading level={2}>Your places</Heading>
                <Spacer />
                <Button variant="secondary" onClick={() => onAction("exportPlaces")}>Export</Button>
              </Row>
              {data.noPlaces && (
                <Text>Find a place or click the map, then save it to a list.</Text>
              )}
              {data.places.map((item, i) => (
                <Row key={i} variant="saved-place">
                  <Button variant="result" onClick={() => onAction("savedPlace")} item={item.id}>
                    <Stack>
                      <Text>{item.name}</Text>
                      <Text variant="caption">{item.list}</Text>
                    </Stack>
                  </Button>
                  <Button variant="icon" onClick={() => onAction("deletePlace")} item={item.id} aria-label="Delete saved place">×</Button>
                </Row>
              ))}
            </>
          )}
          {data.routing && (
            <>
              <Heading level={2}>Directions</Heading>
              <Select name="mode" options={data.modes} aria-label="Travel mode" />
              <Text variant="caption">Pick places on the map or search, then use Set start / Add stop.</Text>
              {data.stops.map((item, i) => (
                <Row key={i}>
                  <Text>{item.label}</Text>
                  <Spacer />
                  <Button variant="icon" onClick={() => onAction("removeStop")} item={item.key} aria-label="Remove stop">×</Button>
                </Row>
              ))}
              <Row>
                <Button onClick={() => onAction("route")} disabled={data.notEnoughStops}>Find routes</Button>
                <Button variant="secondary" onClick={() => onAction("clearStops")}>Clear</Button>
              </Row>
              {data.hasRoutes && (
                <>
                  {data.routes.map((item, i) => (
                    <Button key={i} variant={item.variant} onClick={() => onAction("chooseRoute")} item={item.key}><Text>{item.summary}</Text></Button>
                  ))}
                  <Row>
                    <Button variant="secondary" onClick={() => onAction("exportRoute")}>Export GPX</Button>
                    <Button variant="secondary" onClick={() => onAction("navigate")}>{data.navigationLabel}</Button>
                  </Row>
                  <Text variant="caption">Routes are estimates from OpenStreetMap. No live traffic or road-closure feed.</Text>
                  {data.steps.map((item, i) => (
                    <Row key={i} variant="step">
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
              <Row>
                <Heading level={2}>{data.selected.name}</Heading>
                <Spacer />
                <Button variant="icon" onClick={() => onAction("deselect")} aria-label="Close place details">×</Button>
              </Row>
              <Text>{data.selected.address}</Text>
              <Text variant="caption">{data.selected.coordinates}</Text>
              <Row>
                <Button variant="secondary" onClick={() => onAction("start")}>Set start</Button>
                <Button variant="secondary" onClick={() => onAction("stop")}>Add stop</Button>
              </Row>
              <Field label="List"><Input name="list" placeholder="Favourites" /></Field>
              <Field label="Notes"><Input name="notes" multiline={true} /></Field>
              <Row>
                <Button onClick={() => onAction("save")}>Save place</Button>
                <Button variant="secondary" onClick={() => onAction("share")}>Share</Button>
              </Row>
            </Card>
          )}
          <Spacer />
          <Link href="https://www.openstreetmap.org/fixthemap" external={true}>Improve this map</Link>
        </Stack>
        <Section id="map" />
      </Row>
      <Row variant="map-footer">
        <Button variant="secondary" onClick={() => onAction("locate")}>My location</Button>
        <Button variant="secondary" onClick={() => onAction("measure")}>{data.measureLabel}</Button>
        <Text variant="caption">{data.measurement}</Text>
        <Spacer />
        <Link href="https://www.openstreetmap.org/copyright" external={true}>© OpenStreetMap contributors</Link>
      </Row>
    </Page>
  )
}
