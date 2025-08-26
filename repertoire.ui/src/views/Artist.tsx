import { Divider, Flex, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetArtistQuery } from '../state/api/artistsApi.ts'
import ArtistLoader from '../components/artist/loader/ArtistLoader.tsx'
import { useGetAlbumsQuery } from '../state/api/albumsApi.ts'
import { useGetInfiniteSongsInfiniteQuery } from '../state/api/songsApi.ts'
import { useEffect } from 'react'
import artistSongsOrders from '../data/artist/artistSongsOrders.ts'
import artistAlbumsOrders from '../data/artist/artistAlbumsOrders.ts'
import ArtistAlbumsWidget from '../components/artist/widgets/ArtistAlbumsWidget.tsx'
import ArtistSongsWidget from '../components/artist/widgets/ArtistSongsWidget.tsx'
import ArtistHeader from '../components/artist/ArtistHeader.tsx'
import useDynamicDocumentTitle from '../hooks/useDynamicDocumentTitle.ts'
import BandMembersWidget from '../components/artist/widgets/BandMembersWidget.tsx'
import LocalStorageKeys from '../types/enums/keys/LocalStorageKeys.ts'
import useOrderBy from '../hooks/api/useOrderBy.ts'
import useLocalStorage from '../hooks/useLocalStorage.ts'
import useSearchBy from '../hooks/api/useSearchBy.ts'
import AlbumProperty from '../types/enums/properties/AlbumProperty.ts'
import FilterOperator from '../types/enums/FilterOperator.ts'
import SongProperty from '../types/enums/properties/SongProperty.ts'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import Song from '../types/models/Song.ts'

function Artist() {
  const params = useParams()
  const setDocumentTitle = useDynamicDocumentTitle()
  const artistId = params['id'] ?? ''
  const isUnknownArtist = artistId === 'unknown'

  const {
    data: artist,
    isLoading,
    isFetching
  } = useGetArtistQuery(artistId, { skip: isUnknownArtist })

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
    data: dataSongs,
    isLoading: isSongsLoading,
    isFetching: isSongsFetching,
    isFetchingNextPage: isSongsFetchingNextPage,
    fetchNextPage: songsFetchNextPage
  } = useGetInfiniteSongsInfiniteQuery({
    orderBy: songsOrderBy,
    searchBy: songsSearchBy,
    pageSize: 25
  })
  const songs: WithTotalCountResponse<Song> = {
    models: dataSongs?.pages.flatMap((x) => x.models ?? []),
    totalCount: dataSongs?.pages[0].totalCount
  }

  if (isLoading || (!artist && !isUnknownArtist)) return <ArtistLoader />

  return (
    <Stack h={'100%'} px={'xl'} gap={'16px'}>
      <ArtistHeader
        artist={artist}
        albumsTotalCount={albums?.totalCount}
        songsTotalCount={songs?.totalCount}
        isUnknownArtist={isUnknownArtist}
      />

      <Divider />

      <Flex direction={{ base: 'column', md: 'row' }} mih={360} pb={'lg'} gap={'md'}>
        <Stack
          w={{ base: '100%', md: '55%' }}
          h={{ base: 'unset', md: '100%' }}
          mah={{ base: !isUnknownArtist && artist.isBand ? '60vh' : '37vh', md: 'unset' }}
        >
          {!isUnknownArtist && artist.isBand && (
            <Stack>
              <BandMembersWidget
                bandMembers={artist.bandMembers}
                artistId={artistId}
                isFetching={isFetching}
              />
            </Stack>
          )}

          <ArtistAlbumsWidget
            albums={albums}
            isLoading={isAlbumsLoading}
            isFetching={isAlbumsFetching}
            isUnknownArtist={isUnknownArtist}
            order={albumsOrder}
            setOrder={setAlbumsOrder}
            artistId={artist?.id}
          />
        </Stack>

        <Stack flex={1} h={'100%'} pb={{ base: 'lg', md: 0 }}>
          <ArtistSongsWidget
            songs={songs}
            isUnknownArtist={isUnknownArtist}
            order={songsOrder}
            setOrder={setSongsOrder}
            artistId={artist?.id}
            isLoading={isSongsLoading}
            isFetching={isSongsFetching}
            isFetchingNextPage={isSongsFetchingNextPage}
            fetchNextPage={songsFetchNextPage}
          />
        </Stack>
      </Flex>
    </Stack>
  )
}

export default Artist
