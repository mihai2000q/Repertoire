import { Divider, Grid, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import SongLoader from '../components/song/SongLoader.tsx'
import { useGetSongQuery } from '../state/api/songsApi.ts'
import SongSectionsCard from '../components/song/panels/SongSectionsCard.tsx'
import SongInformationCard from '../components/song/panels/SongInformationCard.tsx'
import SongLinksCard from '../components/song/panels/SongLinksCard.tsx'
import SongOverallCard from '../components/song/panels/SongOverallCard.tsx'
import SongDescriptionCard from '../components/song/panels/SongDescriptionCard.tsx'
import SongHeaderCard from '../components/song/panels/SongHeaderCard.tsx'
import useDynamicDocumentTitle from '../hooks/useDynamicDocumentTitle.ts'
import { useEffect } from 'react'

function Song() {
  const params = useParams()
  const setDocumentTitle = useDynamicDocumentTitle()
  const songId = params['id'] ?? ''

  const { data: song, isLoading } = useGetSongQuery(songId)

  useEffect(() => {
    if (song) setDocumentTitle(song.title)
  }, [song])

  if (isLoading) return <SongLoader />

  return (
    <Stack px={'xl'}>
      <SongHeaderCard song={song} />

      <Divider />

      <Grid align="start" mb={'lg'}>
        <Grid.Col span={{ sm: 12, md: 4.5 }}>
          <Stack>
            <SongInformationCard song={song} />

            <SongOverallCard song={song} />

            <SongLinksCard song={song} />
          </Stack>
        </Grid.Col>

        <Grid.Col span={{ sm: 12, md: 7.5 }}>
          <Stack>
            <SongDescriptionCard song={song} />

            <SongSectionsCard
              songId={songId}
              settings={song.settings}
              sections={song.sections}
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
