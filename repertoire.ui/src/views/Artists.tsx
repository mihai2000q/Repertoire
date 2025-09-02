import { useGetArtistsQuery } from '../state/api/artistsApi.ts'
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
import ArtistsLoader from '../components/artists/ArtistsLoader.tsx'
import ArtistCard from '../components/artists/ArtistCard.tsx'
import AddNewArtistModal from '../components/artists/modal/AddNewArtistModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconArrowsSort, IconFilterFilled, IconPlus, IconUserPlus } from '@tabler/icons-react'
import usePaginationInfo from '../hooks/usePaginationInfo.ts'
import UnknownArtistCard from '../components/artists/UnknownArtistCard.tsx'
import useShowUnknownArtist from '../hooks/useShowUnknownArtist.ts'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import useSearchParamsState from '../hooks/useSearchParamsState.ts'
import artistsSearchParamsState from '../state/searchParams/ArtistsSearchParamsState.ts'
import AdvancedOrderMenu from '../components/@ui/menu/AdvancedOrderMenu.tsx'
import useLocalStorage from '../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../types/enums/keys/LocalStorageKeys.ts'
import useOrderBy from '../hooks/api/useOrderBy.ts'
import artistsOrders from '../data/artists/artistsOrders.ts'
import { artistPropertyIcons } from '../data/icons/artistPropertyIcons.tsx'
import useSearchBy from '../hooks/api/useSearchBy.ts'
import ArtistsFilters from '../components/artists/ArtistsFilters.tsx'
import artistsFilters from '../data/artists/artistsFilters.ts'
import useSearchParamFilters from '../hooks/filter/useSearchParamFilters.ts'
import { DragSelectProvider } from '../context/DragSelectContext.tsx'
import { useRef } from 'react'
import ArtistsSelectionDrawer from '../components/artists/ArtistsSelectionDrawer.tsx'
import ArtistsContextMenu from '../components/artists/ArtistsContextMenu.tsx'
import { useMainScroll } from '../context/MainScrollContext.tsx'

function Artists() {
  useFixedDocumentTitle('Artists')
  const [searchParams, setSearchParams] = useSearchParamsState(artistsSearchParamsState)
  const { currentPage, activeFilters } = searchParams

  const [orders, setOrders] = useLocalStorage({
    key: LocalStorageKeys.ArtistsOrders,
    defaultValue: artistsOrders
  })
  const orderBy = useOrderBy(orders)

  const [filters, setFilters] = useSearchParamFilters({
    initialFilters: artistsFilters,
    activeFilters: activeFilters,
    setSearchParams: setSearchParams
  })
  const filtersSize = activeFilters.size
  const searchBy = useSearchBy(filters)

  const pageSize = 40
  const {
    data: artists,
    isLoading,
    isFetching
  } = useGetArtistsQuery({
    pageSize: pageSize,
    currentPage: currentPage,
    orderBy: orderBy,
    searchBy: searchBy
  })

  const showUnknownArtist = useShowUnknownArtist()
  const [filtersOpened, { toggle: toggleFilters, close: closeFilters }] = useDisclosure(false)

  const { startCount, endCount, totalPages } = usePaginationInfo(
    artists?.totalCount + (showUnknownArtist ? 1 : 0),
    pageSize,
    currentPage
  )

  const [openedAddNewArtistModal, { open: openAddNewArtistModal, close: closeAddNewArtistModal }] =
    useDisclosure(false)

  const artistsRef = useRef<HTMLDivElement>()
  const { ref: mainScrollRef } = useMainScroll()

  function handleCurrentPageChange(p: number) {
    mainScrollRef.current.scrollTo({ top: 0, behavior: 'instant' })
    if (currentPage === p) return
    setSearchParams({ ...searchParams, currentPage: p })
  }

  return (
    <Stack h={'100%'} gap={0}>
      <Group gap={'xxs'} pb={'xs'} px={'xl'}>
        <Title order={3} fw={800} fz={'max(2.5vw, 32px)'}>
          Artists
        </Title>
        <ActionIcon
          aria-label={'new-artist'}
          variant={'grey'}
          size={'lg'}
          onClick={openAddNewArtistModal}
        >
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <AdvancedOrderMenu
          orders={orders}
          setOrders={setOrders}
          propertyIcons={artistPropertyIcons}
        >
          <ActionIcon
            aria-label={'order-artists'}
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
            aria-label={'filter-artists'}
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
        <Text lh={'xxs'} px={'xl'}>
          {startCount} - {endCount} artists out of{' '}
          {(artists?.totalCount ?? 0) + (showUnknownArtist ? 1 : 0)}
        </Text>
      )}

      {artists?.totalCount === 0 && !showUnknownArtist && filtersSize === 0 && (
        <Text p={'xl'}>There are no artists yet. Try to add one</Text>
      )}
      {artists?.totalCount === 0 && !showUnknownArtist && filtersSize > 0 && (
        <Text p={'xl'}>There are no artists with these filter properties</Text>
      )}
      <DragSelectProvider settings={{ area: artistsRef.current }} data={artists}>
        <ArtistsContextMenu>
          <SimpleGrid
            data-testid={'artists-area'}
            ref={artistsRef}
            cols={{ base: 3, xs: 4, sm: 3, betweenSmMd: 4, md: 5, lg: 6, xl: 7, xxl: 8 }}
            verticalSpacing={{ base: 'lg', md: 'xl' }}
            spacing={{ base: 'lg', md: 'xl' }}
            pt={'lg'}
            pb={'xs'}
            px={'xl'}
          >
            {(isLoading || !artists) && <ArtistsLoader />}
            {artists?.models.map((artist) => (
              <ArtistCard key={artist.id} artist={artist} />
            ))}
            {!isFetching && showUnknownArtist && currentPage == totalPages && <UnknownArtistCard />}
            {!isFetching &&
              ((artists?.totalCount > 0 && currentPage == totalPages) ||
                (artists?.totalCount === 0 && showUnknownArtist)) && (
                <Card
                  variant={'add-new'}
                  aria-label={'new-artist-card'}
                  radius={'50%'}
                  onClick={openAddNewArtistModal}
                  style={{ aspectRatio: 1 }}
                >
                  <Center h={'100%'}>
                    <IconUserPlus size={'100%'} style={{ padding: '29%' }} />
                  </Center>
                </Card>
              )}
          </SimpleGrid>
        </ArtistsContextMenu>
        <ArtistsSelectionDrawer />
      </DragSelectProvider>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'md'}>
        <Pagination
          data-testid={'artists-pagination'}
          value={currentPage}
          onChange={handleCurrentPageChange}
          total={totalPages}
          disabled={isFetching}
        />
      </Box>

      <ArtistsFilters
        opened={filtersOpened}
        onClose={closeFilters}
        filters={filters}
        setFilters={setFilters}
        isArtistsLoading={isFetching}
      />
      <AddNewArtistModal opened={openedAddNewArtistModal} onClose={closeAddNewArtistModal} />
    </Stack>
  )
}

export default Artists
