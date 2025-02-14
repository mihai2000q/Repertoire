import { useGetArtistsQuery } from '../state/api/artistsApi.ts'
import {
  ActionIcon,
  Box,
  Card,
  Center,
  Group,
  Loader,
  Pagination,
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

function Artists() {
  useFixedDocumentTitle('Artists')
  const [searchParams, setSearchParams] = useSearchParamsState(artistsSearchParamsState)
  const { currentPage } = searchParams

  const pageSize = 40
  const {
    data: artists,
    isLoading,
    isFetching
  } = useGetArtistsQuery({
    pageSize: pageSize,
    currentPage: currentPage,
    orderBy: ['created_at DESC']
  })

  const showUnknownArtist = useShowUnknownArtist()

  const { startCount, endCount, totalPages } = usePaginationInfo(
    artists?.totalCount + (showUnknownArtist ? 1 : 0),
    pageSize,
    currentPage
  )

  const [openedAddNewArtistModal, { open: openAddNewArtistModal, close: closeAddNewArtistModal }] =
    useDisclosure(false)

  const handleCurrentPageChange = (p: number) => {
    setSearchParams({ ...searchParams, currentPage: p })
  }

  return (
    <Stack h={'100%'} gap={'xs'} px={'xl'}>
      <Group gap={'xxs'}>
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
        <ActionIcon aria-label={'order-artists'} variant={'grey'} size={'lg'} disabled={isLoading}>
          <IconArrowsSort size={17} />
        </ActionIcon>
        <ActionIcon aria-label={'filter-artists'} variant={'grey'} size={'lg'} disabled={isLoading}>
          <IconFilterFilled size={17} />
        </ActionIcon>
      </Group>
      {!isLoading && (
        <Text inline mb={'xs'}>
          {startCount} - {endCount} artists out of{' '}
          {artists?.totalCount + (showUnknownArtist ? 1 : 0)}
        </Text>
      )}

      {artists?.totalCount === 0 && !showUnknownArtist && (
        <Text mt={'sm'}>There are no artists yet. Try to add one</Text>
      )}
      <Group gap={'xl'} align={'start'}>
        {isLoading && <ArtistsLoader />}
        {artists?.models.map((artist) => <ArtistCard key={artist.id} artist={artist} />)}
        {showUnknownArtist && currentPage == totalPages && <UnknownArtistCard />}
        {((artists?.totalCount > 0 && currentPage == totalPages) ||
          (artists?.totalCount === 0 && showUnknownArtist)) && (
          <Card
            aria-label={'new-artist-card'}
            w={125}
            h={125}
            radius={'50%'}
            onClick={openAddNewArtistModal}
            variant={'add-new'}
          >
            <Center h={'100%'}>
              <IconUserPlus size={40} />
            </Center>
          </Card>
        )}
      </Group>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'md'}>
        {!isFetching ? (
          <Pagination
            data-testid={'artists-pagination'}
            value={currentPage}
            onChange={handleCurrentPageChange}
            total={totalPages}
          />
        ) : (
          <Loader size={25} />
        )}
      </Box>

      <AddNewArtistModal opened={openedAddNewArtistModal} onClose={closeAddNewArtistModal} />
    </Stack>
  )
}

export default Artists
