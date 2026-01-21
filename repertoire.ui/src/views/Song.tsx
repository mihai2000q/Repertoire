import { Divider, Grid, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import SongLoader from '../components/song/SongLoader.tsx'
import { useGetSongQuery } from '../state/api/songsApi.ts'
import SongSectionsWidget from '../components/song/widgets/SongSectionsWidget.tsx'
import SongInformationWidget from '../components/song/widgets/SongInformationWidget.tsx'
import SongLinksWidget from '../components/song/widgets/SongLinksWidget.tsx'
import SongOverallWidget from '../components/song/widgets/SongOverallWidget.tsx'
import SongDescriptionWidget from '../components/song/widgets/SongDescriptionWidget.tsx'
import SongHeader from '../components/song/SongHeader.tsx'
import useDynamicDocumentTitle from '../hooks/useDynamicDocumentTitle.ts'
import { useEffect } from 'react'

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

            <SongSectionsWidget
              songId={songId}
              settings={song.settings}
              sections={song.sections}
              isFetching={isFetching}
              bandMembers={song.artist?.isBand === false ? undefined : song.artist?.bandMembers}
              isArtistBand={song.artist?.isBand}
            />
          </Stack>
        </Grid.Col>
      </Grid>
    </Stack>
  )
}

export default Song
