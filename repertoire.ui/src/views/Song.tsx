import { Divider, Grid, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import SongLoader from '../components/song/SongLoader.tsx'
import { useGetSongQuery } from '../state/songsApi.ts'
import SongSections from '../components/song/SongSections.tsx'
import SongInformationCard from '../components/song/panels/SongInformationCard.tsx'
import SongLinksCard from '../components/song/panels/SongLinksCard.tsx'
import SongOverallCard from '../components/song/panels/SongOverallCard.tsx'
import SongDescriptionCard from '../components/song/panels/SongDescriptionCard.tsx'
import SongHeaderCard from '../components/song/panels/SongHeaderCard.tsx'

function Song() {
  const params = useParams()
  const songId = params['id'] ?? ''

  const { data: song, isLoading } = useGetSongQuery(songId)

  if (isLoading) return <SongLoader />

  return (
    <Stack>
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

            <SongSections songId={songId} sections={song.sections} />
          </Stack>
        </Grid.Col>
      </Grid>
    </Stack>
  )
}

export default Song
