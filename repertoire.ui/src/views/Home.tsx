import { memo, ReactElement } from 'react'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import { SimpleGrid, Stack } from '@mantine/core'
import HomeTop from '../components/home/HomeTop.tsx'
import HomeRecentlyPlayedSongs from '../components/home/HomeRecentlyPlayedSongs.tsx'
import HomeGenres from '../components/home/HomeGenres.tsx'
import HomeRecentPlaylists from '../components/home/HomeRecentPlaylists.tsx'
import HomeRecentArtists from '../components/home/HomeRecentArtists.tsx'
import { useElementSize } from '@mantine/hooks'
import useTitleBarHeight from '../hooks/useTitleBarHeight.ts'
import useTopbarHeight from '../hooks/useTopbarHeight.ts'

function Home(): ReactElement {
  useFixedDocumentTitle('Home')
  const stackGap = '16px'
  const { ref: topRef, height: topHeight } = useElementSize()

  return (
    <Stack h={'100%'} gap={stackGap}>
      <HomeTop ref={topRef} />
      <HomeContent stackGap={stackGap} topHeight={topHeight} />
    </Stack>
  )
}

const HomeContent = memo(function HomeContent({
  stackGap,
  topHeight
}: {
  stackGap: string
  topHeight: number
}) {
  const titleBarHeight = useTitleBarHeight()
  const topbarHeight = useTopbarHeight()

  return (
    <SimpleGrid
      px={'xl'}
      cols={{ base: 1, md: 2, lg: 3 }}
      h={`calc(100vh - ${topHeight}px - ${topbarHeight} - ${titleBarHeight} - ${stackGap})`}
      pb={'lg'}
      mih={300}
    >
      <HomeGenres visibleFrom={'md'} />
      <HomeRecentlyPlayedSongs />
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
  )
})

export default Home
