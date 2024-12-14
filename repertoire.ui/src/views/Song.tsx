import { Divider, Group, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import SongLoader from '../components/song/SongLoader.tsx'
import { useGetSongQuery } from '../state/songsApi.ts'
import SongSections from '../components/song/SongSections.tsx'
import SongInformationCard from '../components/song/panels/SongInformationCard.tsx'
import SongLinksCard from '../components/song/panels/SongLinksCard.tsx'
import SongOverallCard from '../components/song/panels/SongsOverallCard.tsx'
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

      <Group align="start" mb={'lg'}>
        <Stack flex={1}>
          <SongInformationCard song={song} />

          <SongOverallCard song={song} />

          <SongLinksCard song={song} />
        </Stack>

        <Stack flex={1.75}>
          <SongDescriptionCard song={song} />

          <SongSections songId={songId} sections={song.sections} />
        </Stack>
      </Group>
    </Stack>
  )
}

export default Song
