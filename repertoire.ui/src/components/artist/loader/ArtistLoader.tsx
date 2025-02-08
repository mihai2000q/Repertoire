import { Divider, Grid, Group, Skeleton, Stack } from '@mantine/core'
import ArtistAlbumsLoader from './ArtistAlbumsLoader.tsx'
import ArtistSongsLoader from './ArtistSongsLoader.tsx'
import BandMembersLoader from "./BandMembersLoader.tsx";

function ArtistLoader() {
  return (
    <Stack data-testid={'artist-loader'}>
      <Group align={'start'}>
        <Skeleton radius={'50%'} w={125} h={125} />
        <Stack gap={4} pt={'16px'}>
          <Skeleton w={80} h={15} />
          <Skeleton w={200} h={35} my={6} />
          <Group gap={4}>
            <Skeleton w={65} h={18} />
            <Skeleton w={55} h={18} />
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
