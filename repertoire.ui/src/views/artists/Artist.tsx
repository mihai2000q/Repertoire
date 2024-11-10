import { Avatar, Group, Stack, Title } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetArtistQuery } from '../../state/artistsApi.ts'
import artistPlaceholder from '../../assets/user-placeholder.jpg'
import ArtistLoader from '../../components/artists/loader/ArtistLoader.tsx'

function Artist() {
  const params = useParams()
  const artistId = params['id']

  const { data: artist, isLoading } = useGetArtistQuery(artistId, { skip: !artistId })

  if (isLoading) return <ArtistLoader />

  return (
    <Stack>
      <Group>
        <Avatar src={artist.imageUrl ? artist.imageUrl : artistPlaceholder} size={125} />
        <Title order={3} fw={700}>
          {artist.name}
        </Title>
      </Group>
    </Stack>
  )
}

export default Artist
