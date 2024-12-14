import { Divider, Grid, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetArtistQuery } from '../state/artistsApi.ts'
import ArtistLoader from '../components/artist/loader/ArtistLoader.tsx'
import { useGetAlbumsQuery } from '../state/albumsApi.ts'
import { useGetSongsQuery } from '../state/songsApi.ts'
import { useState } from 'react'
import Order from '../types/Order.ts'
import artistSongsOrders from '../data/artist/artistSongsOrders.ts'
import artistAlbumsOrders from '../data/artist/artistAlbumsOrders.ts'
import ArtistAlbumsCard from '../components/artist/panels/ArtistAlbumsCard.tsx'
import ArtistSongsCard from '../components/artist/panels/ArtistSongsCard.tsx'
import ArtistHeaderCard from '../components/artist/panels/ArtistHeaderCard.tsx'

function Artist() {
  const params = useParams()
  const artistId = params['id'] ?? ''
  const isUnknownArtist = artistId === 'unknown'

  const { data: artist, isLoading } = useGetArtistQuery(artistId, { skip: isUnknownArtist })

  const [albumsOrder, setAlbumsOrder] = useState<Order>(artistAlbumsOrders[0])
  const [songsOrder, setSongsOrder] = useState<Order>(artistSongsOrders[0])

  const {
    data: albums,
    isLoading: isAlbumsLoading,
    isFetching: isAlbumsFetching
  } = useGetAlbumsQuery({
    orderBy: [albumsOrder.value],
    searchBy: [isUnknownArtist ? 'artist_id IS NULL' : `artist_id = '${artistId}'`]
  })
  const {
    data: songs,
    isLoading: isSongsLoading,
    isFetching: isSongsFetching
  } = useGetSongsQuery({
    orderBy: [songsOrder.value],
    searchBy: [isUnknownArtist ? 'songs.artist_id IS NULL' : `songs.artist_id = '${artistId}'`]
  })

  if (isLoading) return <ArtistLoader />

  return (
    <Stack>
      <ArtistHeaderCard
        artist={artist}
        albumsTotalCount={albums?.totalCount}
        songsTotalCount={songs?.totalCount}
        isUnknownArtist={isUnknownArtist}
      />

      <Divider />

      <Grid align={'flex-start'}>
        <Grid.Col span={{ sm: 12, md: 6.5 }}>
          <ArtistAlbumsCard
            albums={albums}
            isLoading={isAlbumsLoading}
            isFetching={isAlbumsFetching}
            isUnknownArtist={isUnknownArtist}
            order={albumsOrder}
            setOrder={setAlbumsOrder}
            artistId={artist?.id}
          />
        </Grid.Col>

        <Grid.Col span={{ sm: 12, md: 5.5 }}>
          <ArtistSongsCard
            songs={songs}
            isLoading={isSongsLoading}
            isFetching={isSongsFetching}
            isUnknownArtist={isUnknownArtist}
            order={songsOrder}
            setOrder={setSongsOrder}
            artistId={artist?.id}
          />
        </Grid.Col>
      </Grid>
    </Stack>
  )
}

export default Artist
