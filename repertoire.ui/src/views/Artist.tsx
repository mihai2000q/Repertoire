import { Divider, Grid, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetArtistQuery } from '../state/api/artistsApi.ts'
import ArtistLoader from '../components/artist/loader/ArtistLoader.tsx'
import { useGetAlbumsQuery } from '../state/api/albumsApi.ts'
import { useGetSongsQuery } from '../state/api/songsApi.ts'
import { useEffect } from 'react'
import artistSongsOrders from '../data/artist/artistSongsOrders.ts'
import artistAlbumsOrders from '../data/artist/artistAlbumsOrders.ts'
import ArtistAlbumsCard from '../components/artist/panels/ArtistAlbumsCard.tsx'
import ArtistSongsCard from '../components/artist/panels/ArtistSongsCard.tsx'
import ArtistHeaderCard from '../components/artist/panels/ArtistHeaderCard.tsx'
import useDynamicDocumentTitle from '../hooks/useDynamicDocumentTitle.ts'
import BandMembersCard from '../components/artist/panels/BandMembersCard.tsx'
import LocalStorageKeys from '../types/enums/LocalStorageKeys.ts'
import useOrderBy from '../hooks/api/useOrderBy.ts'
import useLocalStorage from '../hooks/useLocalStorage.ts'
import useSearchBy from '../hooks/api/useSearchBy.ts'
import AlbumProperty from '../types/enums/AlbumProperty.ts'
import FilterOperator from '../types/enums/FilterOperator.ts'
import SongProperty from '../types/enums/SongProperty.ts'

function Artist() {
  const params = useParams()
  const setDocumentTitle = useDynamicDocumentTitle()
  const artistId = params['id'] ?? ''
  const isUnknownArtist = artistId === 'unknown'

  const { data: artist, isLoading } = useGetArtistQuery(artistId, { skip: isUnknownArtist })

  useEffect(() => {
    if (isUnknownArtist) setDocumentTitle('Unknown Artist')
    else if (artist) setDocumentTitle(artist.name)
  }, [artist, isUnknownArtist])

  const [albumsOrder, setAlbumsOrder] = useLocalStorage({
    key: LocalStorageKeys.ArtistAlbumsOrder,
    defaultValue: artistAlbumsOrders[0]
  })
  const albumsOrderBy = useOrderBy([albumsOrder])
  const albumsSearchBy = useSearchBy(
    isUnknownArtist
      ? [{ property: AlbumProperty.ArtistId, operator: FilterOperator.IsNull }]
      : [{ property: AlbumProperty.ArtistId, operator: FilterOperator.Equal, value: artistId }]
  )

  const [songsOrder, setSongsOrder] = useLocalStorage({
    key: LocalStorageKeys.ArtistSongsOrder,
    defaultValue: artistSongsOrders[0]
  })
  const songsOrderBy = useOrderBy([songsOrder])
  const songsSearchBy = useSearchBy(
    isUnknownArtist
      ? [{ property: SongProperty.ArtistId, operator: FilterOperator.IsNull }]
      : [{ property: SongProperty.ArtistId, operator: FilterOperator.Equal, value: artistId }]
  )

  const {
    data: albums,
    isLoading: isAlbumsLoading,
    isFetching: isAlbumsFetching
  } = useGetAlbumsQuery({
    orderBy: albumsOrderBy,
    searchBy: albumsSearchBy
  })
  const {
    data: songs,
    isLoading: isSongsLoading,
    isFetching: isSongsFetching
  } = useGetSongsQuery({
    orderBy: songsOrderBy,
    searchBy: songsSearchBy
  })

  if (isLoading || (!artist && !isUnknownArtist)) return <ArtistLoader />

  return (
    <Stack px={'xl'}>
      <ArtistHeaderCard
        artist={artist}
        albumsTotalCount={albums?.totalCount}
        songsTotalCount={songs?.totalCount}
        isUnknownArtist={isUnknownArtist}
      />

      <Divider />

      <Grid align={'flex-start'}>
        <Grid.Col span={{ sm: 12, md: 6.5 }}>
          <Stack>
            {!isUnknownArtist && artist.isBand && (
              <BandMembersCard bandMembers={artist.bandMembers} artistId={artistId} />
            )}

            <ArtistAlbumsCard
              albums={albums}
              isLoading={isAlbumsLoading}
              isFetching={isAlbumsFetching}
              isUnknownArtist={isUnknownArtist}
              order={albumsOrder}
              setOrder={setAlbumsOrder}
              artistId={artist?.id}
            />
          </Stack>
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
