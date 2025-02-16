import { Divider, Grid, Group, Skeleton, Stack } from '@mantine/core'
import ArtistAlbumsLoader from './ArtistAlbumsLoader.tsx'
import ArtistSongsLoader from './ArtistSongsLoader.tsx'
import BandMembersLoader from "./BandMembersLoader.tsx";

function ArtistLoader() {
  return (
    <Stack px={'xl'} data-testid={'artist-loader'}>
      <Group align={'start'}>
        <Skeleton radius={'50%'} w={125} h={125} />
        <Stack gap={'xxs'} pt={'16px'}>
          <Skeleton w={80} h={15} />
          <Skeleton w={200} h={37} my={6} />
          <Group gap={'xxs'}>
            <Skeleton w={75} h={16} />
            <Skeleton w={55} h={16} />
            <Skeleton w={45} h={16} />
          </Group>
        </Stack>
      </Group>

      <Divider />

      <Grid align={'flex-start'}>
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
