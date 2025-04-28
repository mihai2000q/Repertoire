import { useGetAlbumsQuery } from '../state/api/albumsApi.ts'
import {
  ActionIcon,
  Box,
  Card,
  Center,
  Group,
  Indicator,
  Pagination,
  SimpleGrid,
  Space,
  Stack,
  Text,
  Title
} from '@mantine/core'
import AlbumsLoader from '../components/albums/AlbumsLoader.tsx'
import AlbumCard from '../components/albums/AlbumCard.tsx'
import AddNewAlbumModal from '../components/albums/modal/AddNewAlbumModal.tsx'
import { useDidUpdate, useDisclosure } from '@mantine/hooks'
import { IconArrowsSort, IconDisc, IconFilterFilled, IconPlus } from '@tabler/icons-react'
import usePaginationInfo from '../hooks/usePaginationInfo.ts'
import useShowUnknownAlbum from '../hooks/useShowUnknownAlbum.ts'
import UnknownAlbumCard from '../components/albums/UnknownAlbumCard.tsx'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import useSearchParamsState from '../hooks/useSearchParamsState.ts'
import albumsSearchParamsState from '../state/searchParams/AlbumsSearchParamsState.ts'
import useOrderBy from '../hooks/api/useOrderBy.ts'
import albumsOrders from '../data/albums/albumsOrders.ts'
import AdvancedOrderMenu from '../components/@ui/menu/AdvancedOrderMenu.tsx'
import LocalStorageKeys from '../types/enums/LocalStorageKeys.ts'
import useLocalStorage from '../hooks/useLocalStorage.ts'
import { albumPropertyIcons } from '../data/icons/albumPropertyIcons.tsx'
import useSearchBy from '../hooks/api/useSearchBy.ts'
import useSearchParamFilters from '../hooks/filter/useSearchParamFilters.ts'
import albumsFilters from '../data/albums/albumsFilters.ts'
import AlbumsFilters from '../components/albums/AlbumsFilters.tsx'

function Albums() {
  useFixedDocumentTitle('Albums')
  const [searchParams, setSearchParams] = useSearchParamsState(albumsSearchParamsState)
  const { currentPage, activeFilters } = searchParams

  const [orders, setOrders] = useLocalStorage({
    key: LocalStorageKeys.AlbumsOrders,
    defaultValue: albumsOrders
  })
  const orderBy = useOrderBy(orders)

  const [filters, setFilters] = useSearchParamFilters({
    initialFilters: albumsFilters,
    activeFilters: activeFilters,
    setSearchParams: setSearchParams
  })
  const filtersSize = activeFilters.size
  const searchBy = useSearchBy(filters)
  useDidUpdate(() => handleCurrentPageChange(1), [JSON.stringify([...activeFilters])])

  const pageSize = 40
  const {
    data: albums,
    isLoading,
    isFetching
  } = useGetAlbumsQuery({
    pageSize: pageSize,
    currentPage: currentPage,
    orderBy: orderBy,
    searchBy: searchBy
  })

  const showUnknownAlbum = useShowUnknownAlbum()
  const [filtersOpened, { toggle: toggleFilters, close: closeFilters }] = useDisclosure(false)

  const { startCount, endCount, totalPages } = usePaginationInfo(
    albums?.totalCount + (showUnknownAlbum ? 1 : 0),
    pageSize,
    currentPage
  )

  const [openedAddNewAlbumModal, { open: openAddNewAlbumModal, close: closeAddNewAlbumModal }] =
    useDisclosure(false)

  function handleCurrentPageChange(p: number) {
    setSearchParams({ ...searchParams, currentPage: p })
  }

  return (
    <Stack h={'100%'} gap={'xs'} px={'xl'}>
      <Group gap={'xxs'}>
        <Title order={3} fw={800} fz={'max(2.5vw, 32px)'}>
          Albums
        </Title>
        <ActionIcon
          aria-label={'new-album'}
          variant={'grey'}
          size={'lg'}
          onClick={openAddNewAlbumModal}
        >
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <AdvancedOrderMenu orders={orders} setOrders={setOrders} propertyIcons={albumPropertyIcons}>
          <ActionIcon aria-label={'order-albums'} variant={'grey'} size={'lg'} disabled={isLoading}>
            <IconArrowsSort size={17} />
          </ActionIcon>
        </AdvancedOrderMenu>
        <Indicator
          color={'primary'}
          size={15}
          offset={3}
          label={filtersSize}
          disabled={filtersSize === 0}
          zIndex={2}
        >
          <ActionIcon
            aria-label={'filter-albums'}
            variant={'grey'}
            size={'lg'}
            disabled={isLoading}
            onClick={toggleFilters}
          >
            <IconFilterFilled size={17} />
          </ActionIcon>
        </Indicator>
      </Group>
      {!isLoading && (
        <Text inline mb={'xs'}>
          {startCount} - {endCount} albums out of{' '}
          {(albums?.totalCount ?? 0) + (showUnknownAlbum ? 1 : 0)}
        </Text>
      )}

      <Stack gap={'xs'} flex={1}>
        {albums?.totalCount === 0 && !showUnknownAlbum && filtersSize === 0 && (
          <Text mt={'sm'}>There are no albums yet. Try to add one</Text>
        )}
        {albums?.totalCount === 0 && !showUnknownAlbum && filtersSize > 0 && (
          <Text mt={'sm'}>There are no albums with these filter properties</Text>
        )}
        <SimpleGrid
          cols={{ base: 2, xs: 3, md: 4, lg: 5, xl: 6, xxl: 7 }}
          verticalSpacing={{ base: 'lg', md: 'xl' }}
          spacing={{ base: 'lg', md: 'xl' }}
        >
          {(isLoading || !albums) && <AlbumsLoader />}
          {albums?.models.map((album) => <AlbumCard key={album.id} album={album} />)}
          {!isFetching && showUnknownAlbum && currentPage == totalPages && <UnknownAlbumCard />}
          {!isFetching &&
            ((albums?.totalCount > 0 && currentPage == totalPages) ||
              (albums?.totalCount === 0 && showUnknownAlbum)) && (
              <Card
                variant={'add-new'}
                aria-label={'new-album-card'}
                radius={'lg'}
                onClick={openAddNewAlbumModal}
                style={{ aspectRatio: 1 }}
              >
                <Center h={'100%'}>
                  <IconDisc size={'100%'} style={{ padding: '37%' }} />
                </Center>
              </Card>
            )}
        </SimpleGrid>
      </Stack>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'md'}>
        <Pagination
          data-testid={'albums-pagination'}
          value={currentPage}
          onChange={handleCurrentPageChange}
          total={totalPages}
          disabled={isFetching}
        />
      </Box>

      <AlbumsFilters
        opened={filtersOpened}
        onClose={closeFilters}
        filters={filters}
        setFilters={setFilters}
        isAlbumsLoading={isFetching}
      />

      <AddNewAlbumModal opened={openedAddNewAlbumModal} onClose={closeAddNewAlbumModal} />
    </Stack>
  )
}

export default Albums
