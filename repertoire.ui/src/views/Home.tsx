import { ReactElement } from 'react'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import { SimpleGrid, Stack } from '@mantine/core'
import HomeTop from '../components/home/HomeTop.tsx'
import HomeRecentlyPlayed from '../components/home/widgets/HomeRecentlyPlayed.tsx'
import HomeGenres from '../components/home/widgets/HomeGenres.tsx'
import HomeRecentPlaylists from '../components/home/widgets/HomeRecentPlaylists.tsx'
import HomeRecentArtists from '../components/home/widgets/HomeRecentArtists.tsx'

function Home(): ReactElement {
  useFixedDocumentTitle('Home')

  return (
    <Stack h={'100%'} gap={16}>
      <HomeTop />

      <SimpleGrid px={'xl'} cols={{ base: 1, md: 2, lg: 3 }} h={'100%'} pb={'lg'} mih={300}>
        <HomeGenres visibleFrom={'md'} />
        <HomeRecentlyPlayed />
        <Stack visibleFrom={'lg'}>
          <HomeRecentArtists
            sx={{
              '@media(max-height: 850px)': {
                display: 'none'
              }
            }}
          />
          <HomeRecentPlaylists flex={1} />
        </Stack>
      </SimpleGrid>
    </Stack>
  )
}

export default Home
