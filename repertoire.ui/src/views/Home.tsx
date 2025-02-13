import { ReactElement } from 'react'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import { SimpleGrid, Stack } from '@mantine/core'
import HomeTop from '../components/home/HomeTop.tsx'
import HomeRecentlyPlayedSongs from '../components/home/HomeRecentlyPlayedSongs.tsx'
import HomeGenres from '../components/home/HomeGenres.tsx'
import HomePlaylists from '../components/home/HomePlaylists.tsx'
import HomeTopArtists from "../components/home/HomeTopArtists.tsx";

function Home(): ReactElement {
  useFixedDocumentTitle('Home')

  return (
    <Stack h={'100%'}>
      <HomeTop />

      <SimpleGrid px={'xl'} cols={{ base: 1, md: 2, lg: 3 }} h={'100%'} mb={'md'} mih={300}>
        <HomeGenres visibleFrom={'md'} />
        <HomeRecentlyPlayedSongs />
        <Stack visibleFrom={'lg'}>
          <HomeTopArtists
            sx={{
              '@media(max-height: 850px)': {
                display: 'none'
              }
            }}
          />
          <HomePlaylists flex={1} />
        </Stack>
      </SimpleGrid>
    </Stack>
  )
}

export default Home
