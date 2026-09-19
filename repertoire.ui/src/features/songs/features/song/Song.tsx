import { Divider, Grid, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import SongLoader from './components/SongLoader.tsx'
import { useGetSongQuery } from '../../../../state/api/songsApi.ts'
import SongInformationWidget from './components/widgets/SongInformationWidget.tsx'
import SongLinksWidget from './components/widgets/SongLinksWidget.tsx'
import SongOverallWidget from './components/widgets/SongOverallWidget.tsx'
import SongDescriptionWidget from './components/widgets/SongDescriptionWidget.tsx'
import SongHeader from './components/SongHeader.tsx'
import useDynamicDocumentTitle from '../../../../hooks/useDynamicDocumentTitle.ts'
import { useEffect } from 'react'
import SongOutlineWidget from './features/outline/SongOutlineWidget.tsx'
import { SongProvider } from './context/SongContext.tsx'
import { SongOutlineProvider } from './features/outline/context/SongOutlineContext.tsx'

function Song() {
  const params = useParams()
  const setDocumentTitle = useDynamicDocumentTitle()
  const songId = params['id'] ?? ''

  const { data: song, isLoading, isFetching } = useGetSongQuery(songId)

  useEffect(() => {
    if (song) setDocumentTitle(song.title)
  }, [song])

  if (isLoading || !song) return <SongLoader />

  return (
    <SongProvider song={song}>
      <Stack px={'xl'}>
        <SongHeader song={song} />

        <Divider />

        <Grid align="start" mb={'lg'}>
          <Grid.Col span={{ sm: 12, md: 4.5 }}>
            <Stack>
              <SongInformationWidget song={song} />

              <SongOverallWidget song={song} />

              <SongLinksWidget song={song} />
            </Stack>
          </Grid.Col>

          <Grid.Col span={{ sm: 12, md: 7.5 }}>
            <Stack>
              <SongDescriptionWidget song={song} />

              <SongOutlineProvider>
                <SongOutlineWidget isSongFetching={isFetching} />
              </SongOutlineProvider>
            </Stack>
          </Grid.Col>
        </Grid>
      </Stack>
    </SongProvider>
  )
}

export default Song
