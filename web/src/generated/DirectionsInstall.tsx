// Generated from screens/DirectionsInstall.aui; edit the AUI source.
import { Alert, Button, Card, Field, Grid, Heading, Icon, Input, Page, Row, Select, Spinner, Stack, Text } from "@/components/ui"

export function DirectionsInstall({ data, onAction }: { data: any; onAction: (name: string) => void }) {
  return (
    <Page>
      <Card variant="directions-download">
        <Row>
          <Icon name="Route" size={28} />
          <Heading level={2}>Local directions</Heading>
        </Row>
        <Text>Drive, walk and cycle using OpenStreetMap data stored on this server. No separate routing service is needed.</Text>
        {data.installed && (
          <>
            <Text variant="route-ready">{data.installedLabel}</Text>
            <Text variant="caption">{data.installedSize}</Text>
          </>
        )}
        {data.error && (
          <Alert tone="error">{data.error}</Alert>
        )}
        {data.owner ? (
          <>
            {data.canPlan && (
              <>
                <Field label="Directions coverage"><Select name="coverage" options={data.coverages} /></Field>
                {data.area && (
                  <>
                    <Grid columns={2}>
                      <Field label="Centre latitude"><Input name="lat" type="number" min={-80} max={80} step={0.0001} /></Field>
                      <Field label="Centre longitude"><Input name="lon" type="number" min={-180} max={180} step={0.0001} /></Field>
                    </Grid>
                    <Field label="Radius in kilometres (1–500)"><Input name="radius" type="number" min={1} max={500} /></Field>
                  </>
                )}
                <Button variant="soft-pill" onClick={() => onAction("check")}>Check directions download</Button>
              </>
            )}
            {data.checking && (
              <Row>
                <Spinner />
                <Text>Checking download size...</Text>
              </Row>
            )}
            {data.plan && (
              <>
                <Text>{data.downloadLabel}</Text>
                <Text variant="caption">{data.spaceLabel}</Text>
                <Text variant="caption">The download includes the local engine and road graph. Five-degree data tiles can extend beyond your chosen area.</Text>
                <Button variant="pill" onClick={() => onAction("install")}>Install local directions</Button>
              </>
            )}
            {data.working ? (
              <>
                <Row>
                  <Spinner />
                  <Text>{data.progressMessage}</Text>
                </Row>
                <Text variant="caption">{data.progressBytes}</Text>
                <Button variant="soft-pill" onClick={() => onAction("cancel")} cancel={true}>Cancel directions download</Button>
              </>
            ) : (
              <Stack>
                {data.progressMessage && (
                  <Text variant="caption">{data.progressMessage}</Text>
                )}
              </Stack>
            )}
          </>
        ) : (
          <Stack>
            {data.needsInstall && (
              <Text>The Workspace owner can install directions here for everyone on this server.</Text>
            )}
          </Stack>
        )}
        <Text variant="caption">Routes include estimates, alternatives and turn instructions. Live traffic, public transport, road closures and charging stops are not included.</Text>
      </Card>
    </Page>
  )
}
