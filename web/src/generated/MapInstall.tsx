// Generated from screens/MapInstall.aui; edit the AUI source.
import { Button, Card, Checkbox, Field, Grid, Heading, Icon, Input, Link, Page, Row, Select, Spinner, Stack, Text } from "@/components/ui"

export function MapInstall({ data, onAction }: { data: any; onAction: (name: string) => void }) {
  return (
    <Page>
      <Stack variant="install-page">
        <Row variant="install-heading">
          <Icon name="MapPinned" size={32} />
          <Stack>
            <Heading level={1}>Your map, your server</Heading>
            <Text variant="caption">Choose how much of the world to keep locally.</Text>
          </Stack>
        </Row>
        {data.installed && (
          <Card>
            <Heading level={2}>Installed map</Heading>
            <Text>{data.installedLabel}</Text>
            <Text variant="caption">{data.installedSize}</Text>
            <Button variant="pill" onClick={() => onAction("openMap")}>Open map</Button>
          </Card>
        )}
        <Card>
          <Heading level={2}>Choose your coverage</Heading>
          <Text>Data downloads only when you choose Install. Browsing, place search and nearby results then stay on this server.</Text>
          {data.owner ? (
            <>
              <Field label="Coverage"><Select name="coverage" options={data.coverages} /></Field>
              {data.area && (
                <>
                  <Grid columns={2} variant="install-coordinates">
                    <Field label="Centre latitude"><Input name="lat" type="number" min={-80} max={80} step={0.0001} /></Field>
                    <Field label="Centre longitude"><Input name="lon" type="number" min={-180} max={180} step={0.0001} /></Field>
                  </Grid>
                  <Field label="Radius in kilometres (1–100)"><Input name="radius" type="number" min={1} max={100} /></Field>
                </>
              )}
              <Field label="Map detail"><Select name="depth" options={data.depths} /></Field>
              <Text variant="caption">{data.depthText}</Text>
              {data.postcodeChoice && (
                <Checkbox name="postcodes" label="Include UK postcode search within this coverage" />
              )}
              <Text variant="caption">Place search includes mapped names and roads, larger towns and optional UK postcodes. This does not include every house address. Street detail provides the most nearby places.</Text>
              {data.canPlan && (
                <Button variant="pill" onClick={() => onAction("plan")}>Check download size</Button>
              )}
            </>
          ) : (
            <Text>The Workspace owner chooses and installs coverage for everyone on this server. Ask the owner to open Map downloads.</Text>
          )}
        </Card>
        {data.plan && (
          <Card>
            <Heading level={2}>Ready to download</Heading>
            <Text>{data.planLabel}</Text>
            <Text>{data.downloadLabel}</Text>
            <Text>{data.spaceLabel}</Text>
            <Text variant="caption">The storage allowance includes indexing and room for safe installation. The existing map remains available until the new one is ready.</Text>
            <Button variant="pill" onClick={() => onAction("install")}>Install selected map</Button>
          </Card>
        )}
        {data.progress && (
          <Card variant="install-progress">
            <Row>
              {data.working && (
                <Spinner />
              )}
              <Heading level={2}>{data.phaseLabel}</Heading>
            </Row>
            <Text>{data.progressMessage}</Text>
            <Text variant="caption">{data.progressBytes}</Text>
            {data.cancelAllowed && (
              <Button variant="soft-pill" onClick={() => onAction("cancel")} cancel={true}>Cancel download</Button>
            )}
          </Card>
        )}
        {data.owner && (
          <Card>
            <Heading level={2}>Local directions (optional)</Heading>
            <Text>Connect an OSRM driving service running on this machine. Directions and Start appear when it is connected.</Text>
            <Field label="Local routing URL"><Input name="router" placeholder="http://127.0.0.1:5000" /></Field>
            <Row>
              <Button variant="soft-pill" onClick={() => onAction("connectRouter")}>Connect routing</Button>
              <Button variant="text" onClick={() => onAction("disconnectRouter")}>Disconnect</Button>
            </Row>
            <Text variant="caption">{data.routerStatus}</Text>
          </Card>
        )}
        <Card>
          <Heading level={2}>Open data, stored locally</Heading>
          <Text>Map data: OpenStreetMap and Protomaps. Town names and UK postcodes: GeoNames. These sources are contacted for your chosen download; there is no Google data or public browsing fallback.</Text>
          <Text variant="caption">UK postcodes contain Royal Mail data. Royal Mail copyright and database right 2022.</Text>
          <Link href="https://www.openstreetmap.org/copyright" external={true}>OpenStreetMap data licence</Link>
          <Link href="https://www.geonames.org/export/" external={true}>GeoNames data licence</Link>
          <Link href="#settings">Settings and encrypted backups</Link>
          <Button variant="text" onClick={() => onAction("logout")}>Sign out</Button>
        </Card>
      </Stack>
    </Page>
  )
}
