import { Divider, Grid, Group, Skeleton, Stack } from '@mantine/core'
import ArtistAlbumsLoader from './ArtistAlbumsLoader.tsx'
import ArtistSongsLoader from './ArtistSongsLoader.tsx'
import BandMembersLoader from "./BandMembersLoader.tsx";

function ArtistLoader() {
  return (
    <Stack px={'xl'} data-testid={'artist-loader'}>
      <Group>
        <Skeleton radius={'50%'} w={'max(12vw, 125px)'} h={'max(12vw, 125px)'} />
        <Stack gap={'xxs'} pt={'16px'}>
          <Skeleton w={80} h={15} />
          <Skeleton w={'max(20vw, 200px)'} h={'max(4vw, 48px)'} my={'xs'} />
          <Group gap={'xxs'}>
            <Skeleton w={75} h={16} />
            <Skeleton w={55} h={16} />
            <Skeleton w={45} h={16} />
          </Group>
        </Stack>
      </Group>

      <Divider />

      <Grid align={'start'}>
        <Grid.Col span={{ sm: 12, md: 6.5 }}>
          <Stack>
            <BandMembersLoader />
            <ArtistAlbumsLoader />
          </Stack>
        </Grid.Col>

        <Grid.Col span={{ sm: 12, md: 5.5 }}>
          <ArtistSongsLoader />
        </Grid.Col>
      </Grid>
    </Stack>
  )
}

export default ArtistLoader
