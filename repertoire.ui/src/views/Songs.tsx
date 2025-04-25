import { ReactElement } from 'react'
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
import { useGetSongsQuery } from '../state/api/songsApi.ts'
import SongCard from '../components/songs/SongCard.tsx'
import { IconArrowsSort, IconFilterFilled, IconMusicPlus, IconPlus } from '@tabler/icons-react'
import { useDidUpdate, useDisclosure } from '@mantine/hooks'
import AddNewSongModal from '../components/songs/modal/AddNewSongModal.tsx'
import SongsLoader from '../components/songs/SongsLoader.tsx'
import usePaginationInfo from '../hooks/usePaginationInfo.ts'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import useSearchParamsState from '../hooks/useSearchParamsState.ts'
import songsSearchParamsState from '../state/searchParams/SongsSearchParamsState.ts'
import useLocalStorage from '../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../types/enums/LocalStorageKeys.ts'
import useOrderBy from '../hooks/api/useOrderBy.ts'
import songsOrders from '../data/songs/songsOrders.ts'
import AdvancedOrderMenu from '../components/@ui/menu/AdvancedOrderMenu.tsx'
import { songPropertyIcons } from '../data/icons/songPropertyIcons.tsx'
import useFilters from '../hooks/filter/useFilters.ts'
import useSearchBy from '../hooks/api/useSearchBy.ts'
import SongsFilters from '../components/songs/SongsFilters.tsx'
import songsFilters from '../data/songs/songsFilters.ts'

function Songs(): ReactElement {
  useFixedDocumentTitle('Songs')
  const [searchParams, setSearchParams] = useSearchParamsState(songsSearchParamsState)
  const { currentPage } = searchParams

  const [orders, setOrders] = useLocalStorage({
    key: LocalStorageKeys.SongsOrders,
    defaultValue: songsOrders
  })
  const orderBy = useOrderBy(orders)

  const [filters, filtersSize, setFilters] = useFilters(songsFilters)
  const searchBy = useSearchBy(filters)
  useDidUpdate(() => handleCurrentPageChange(1), [filters])

  const pageSize = 40
  const {
    data: songs,
    isLoading,
    isFetching
  } = useGetSongsQuery({
    pageSize: pageSize,
    currentPage: currentPage,
    orderBy: orderBy,
    searchBy: searchBy
  })

  const [filtersOpened, { toggle: toggleFilters, close: closeFilters }] = useDisclosure(false)

  const { startCount, endCount, totalPages } = usePaginationInfo(
    songs?.totalCount,
    pageSize,
    currentPage
  )

  const [openedAddNewSongModal, { open: openAddNewSongModal, close: closeAddNewSongModal }] =
    useDisclosure(false)

  const handleCurrentPageChange = (p: number) => {
    setSearchParams({ ...searchParams, currentPage: p })
  }

  return (
    <Stack h={'100%'} gap={'xs'} px={'xl'}>
      <Group gap={'xxs'}>
        <Title order={3} fw={800} fz={'max(2.5vw, 32px)'}>
          Songs
        </Title>
        <ActionIcon
          aria-label={'new-song'}
          variant={'grey'}
          size={'lg'}
          onClick={openAddNewSongModal}
        >
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <AdvancedOrderMenu orders={orders} setOrders={setOrders} propertyIcons={songPropertyIcons}>
          <ActionIcon aria-label={'order-songs'} variant={'grey'} size={'lg'} disabled={isLoading}>
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
            aria-label={'filter-songs'}
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
          {startCount} - {endCount} songs out of {songs?.totalCount ?? 0}
        </Text>
      )}

      {songs?.totalCount === 0 && filtersSize === 0 && (
        <Text mt={'sm'}>There are no songs yet. Try to add one</Text>
      )}
      {songs?.totalCount === 0 && filtersSize > 0 && (
        <Text mt={'sm'}>There are no songs with these filter properties</Text>
      )}
      <SimpleGrid
        cols={{ base: 2, xs: 3, sm: 2, betweenSmMd: 3, betweenMdLg: 4, lg: 5, xl: 6, xxl: 7 }}
        spacing={'lg'}
        verticalSpacing={'lg'}
      >
        {(isLoading || !songs) && <SongsLoader />}
        {songs?.models.map((song) => <SongCard key={song.id} song={song} />)}
        {!isFetching && songs?.totalCount > 0 && currentPage == totalPages && (
          <Card
            variant={'add-new'}
            aria-label={'new-song-card'}
            radius={'lg'}
            onClick={openAddNewSongModal}
          >
            <Center h={'100%'}>
              <IconMusicPlus size={'100%'} style={{ padding: '35%' }} />
            </Center>
          </Card>
        )}
      </SimpleGrid>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'md'}>
        <Pagination
          data-testid={'songs-pagination'}
          value={currentPage}
          onChange={handleCurrentPageChange}
          total={totalPages}
          disabled={isFetching}
        />
      </Box>

      <SongsFilters
        opened={filtersOpened}
        onClose={closeFilters}
        filters={filters}
        setFilters={setFilters}
        isSongsLoading={isLoading}
      />

      <AddNewSongModal opened={openedAddNewSongModal} onClose={closeAddNewSongModal} />
    </Stack>
  )
}

export default Songs
