import { useGetPlaylistsQuery } from '../state/api/playlistsApi.ts'
import usePaginationInfo from '../hooks/usePaginationInfo.ts'
import { useDisclosure } from '@mantine/hooks'
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
import AddNewPlaylistModal from '../components/playlists/modal/AddNewPlaylistModal.tsx'
import { IconArrowsSort, IconFilterFilled, IconPlaylistAdd, IconPlus } from '@tabler/icons-react'
import PlaylistsLoader from '../components/playlists/PlaylistsLoader.tsx'
import PlaylistCard from '../components/playlists/PlaylistCard.tsx'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import useSearchParamsState from '../hooks/useSearchParamsState.ts'
import AdvancedOrderMenu from '../components/@ui/menu/AdvancedOrderMenu.tsx'
import { playlistPropertyIcons } from '../data/icons/playlistPropertyIcons.tsx'
import useLocalStorage from '../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../types/enums/keys/LocalStorageKeys.ts'
import useOrderBy from '../hooks/api/useOrderBy.ts'
import playlistsOrders from '../data/playlists/playlistsOrders.ts'
import PlaylistsFilters from '../components/playlists/PlaylistsFilters.tsx'
import useSearchBy from '../hooks/api/useSearchBy.ts'
import playlistsFilters from '../data/playlists/playlistsFilters.ts'
import useSearchParamFilters from '../hooks/filter/useSearchParamFilters.ts'
import playlistsSearchParamsState from '../state/searchParams/PlaylistsSearchParamsState.ts'
import { useRef } from 'react'
import { useMain } from '../context/MainContext.tsx'
import PlaylistsSelectionDrawer from '../components/playlists/PlaylistsSelectionDrawer.tsx'
import PlaylistsContextMenu from '../components/playlists/PlaylistsContextMenu.tsx'
import { DragSelectProvider } from '../context/DragSelectContext.tsx'

function Playlists() {
  useFixedDocumentTitle('Playlists')
  const [searchParams, setSearchParams] = useSearchParamsState(playlistsSearchParamsState)
  const { currentPage, activeFilters } = searchParams

  const [orders, setOrders] = useLocalStorage({
    key: LocalStorageKeys.PlaylistsOrders,
    defaultValue: playlistsOrders
  })
  const orderBy = useOrderBy(orders)

  const [filters, setFilters] = useSearchParamFilters({
    initialFilters: playlistsFilters,
    activeFilters: activeFilters,
    setSearchParams: setSearchParams
  })
  const filtersSize = activeFilters.size
  const searchBy = useSearchBy(filters)

  const pageSize = 40
  const {
    data: playlists,
    isLoading,
    isFetching
  } = useGetPlaylistsQuery({
    pageSize: pageSize,
    currentPage: currentPage,
    orderBy: orderBy,
    searchBy: searchBy
  })

  const [filtersOpened, { toggle: toggleFilters, close: closeFilters }] = useDisclosure(false)

  const { startCount, endCount, totalPages } = usePaginationInfo(
    playlists?.totalCount,
    pageSize,
    currentPage
  )

  const [
    openedAddNewPlaylistModal,
    { open: openAddNewPlaylistModal, close: closeAddNewPlaylistModal }
  ] = useDisclosure(false)

  const playlistsRef = useRef()
  const { mainScroll } = useMain()

  function handleCurrentPageChange(p: number) {
    mainScroll.ref.current?.scrollTo({ top: 0, behavior: 'instant' })
    if (currentPage === p) return
    setSearchParams({ ...searchParams, currentPage: p })
  }

  return (
    <Stack h={'100%'} gap={0}>
      <Group px={'xl'} gap={'xxs'} pb={'xs'}>
        <Title order={3} fw={800} fz={'max(2.5vw, 32px)'}>
          Playlists
        </Title>
        <ActionIcon
          aria-label={'new-playlist'}
          variant={'grey'}
          size={'lg'}
          onClick={openAddNewPlaylistModal}
        >
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <AdvancedOrderMenu
          orders={orders}
          setOrders={setOrders}
          propertyIcons={playlistPropertyIcons}
        >
          <ActionIcon
            aria-label={'order-playlists'}
            variant={'grey'}
            size={'lg'}
            disabled={isLoading}
          >
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
            aria-label={'filter-playlists'}
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
        <Text px={'xl'} lh={'xxs'}>
          {startCount} - {endCount} playlists out of {playlists?.totalCount ?? 0}
        </Text>
      )}

      {playlists?.totalCount === 0 && filtersSize === 0 && (
        <Text p={'xl'}>There are no playlists yet. Try to add one</Text>
      )}
      {playlists?.totalCount === 0 && filtersSize > 0 && (
        <Text p={'xl'}>There are no playlists with these filters</Text>
      )}
      <DragSelectProvider settings={{ area: playlistsRef.current }} data={playlists}>
        <PlaylistsContextMenu>
          <SimpleGrid
            data-testid={'playlists-area'}
            ref={playlistsRef}
            cols={{ base: 2, xs: 3, md: 4, lg: 5, xl: 6, xxl: 7 }}
            spacing={'lg'}
            verticalSpacing={'lg'}
            pt={'lg'}
            pb={'xs'}
            px={'xl'}
          >
            {(isLoading || !playlists) && <PlaylistsLoader />}
            {playlists?.models.map((playlist) => (
              <PlaylistCard key={playlist.id} playlist={playlist} />
            ))}
            {!isFetching && playlists?.totalCount > 0 && currentPage == totalPages && (
              <Card
                aria-label={'new-playlist-card'}
                variant={'add-new'}
                radius={'lg'}
                onClick={openAddNewPlaylistModal}
                style={{ aspectRatio: 1 }}
              >
                <Center h={'100%'}>
                  <IconPlaylistAdd size={'100%'} style={{ padding: '33%' }} />
                </Center>
              </Card>
            )}
          </SimpleGrid>
        </PlaylistsContextMenu>
        <PlaylistsSelectionDrawer />
      </DragSelectProvider>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'md'}>
        <Pagination
          data-testid={'playlists-pagination'}
          value={currentPage}
          onChange={handleCurrentPageChange}
          total={totalPages}
          disabled={isFetching}
        />
      </Box>

      <PlaylistsFilters
        opened={filtersOpened}
        onClose={closeFilters}
        filters={filters}
        setFilters={setFilters}
        isPlaylistsLoading={isFetching}
      />

      <AddNewPlaylistModal opened={openedAddNewPlaylistModal} onClose={closeAddNewPlaylistModal} />
    </Stack>
  )
}

export default Playlists
