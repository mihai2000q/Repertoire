import { useState } from 'react'
import { useGetArtistsQuery } from '../state/artistsApi.ts'
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

function Artists() {
  const [currentPage, setCurrentPage] = useState(1)
  const { data: artists, isLoading } = useGetArtistsQuery({
    pageSize: 20,
    currentPage: currentPage,
    orderBy: ['created_at DESC']
  })

  const showUnknownArtist = useShowUnknownArtist()

  const { startCount, endCount, totalPages } = usePaginationInfo(
    artists?.totalCount + (showUnknownArtist ? 1 : 0),
    20,
    currentPage
  )

  const [openedAddNewArtistModal, { open: openAddNewArtistModal, close: closeAddNewArtistModal }] =
    useDisclosure(false)

  return (
    <Stack h={'100%'} gap={'xs'}>
      <AddNewArtistModal opened={openedAddNewArtistModal} onClose={closeAddNewArtistModal} />

      <Group gap={4} align={'center'}>
        <Title order={3} fw={800}>
          Artists
        </Title>
        <ActionIcon variant={'grey'} size={'lg'} onClick={openAddNewArtistModal}>
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <ActionIcon variant={'grey'} size={'lg'}>
          <IconArrowsSort size={17} />
        </ActionIcon>
        <ActionIcon variant={'grey'} size={'lg'}>
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
      <Group gap={'xl'}>
        {isLoading && <ArtistsLoader />}
        {artists?.models.map((artist) => <ArtistCard key={artist.id} artist={artist} />)}
        {showUnknownArtist && currentPage == totalPages && <UnknownArtistCard />}
        {((artists?.totalCount > 0 && currentPage == totalPages) || showUnknownArtist) && (
          <Card
            data-testid={'new-artist-card'}
            w={125}
            h={125}
            radius={'50%'}
            onClick={openAddNewArtistModal}
            style={{ alignSelf: 'start' }}
            variant={'add-new'}
          >
            <Center h={'100%'}>
              <IconUserPlus size={40} />
            </Center>
          </Card>
        )}
      </Group>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'xs'}>
        {!isLoading ? (
          <Pagination
            data-testid={'artists-pagination'}
            value={currentPage}
            onChange={setCurrentPage}
            total={totalPages}
          />
        ) : (
          <Loader size={25} />
        )}
      </Box>
    </Stack>
  )
}

export default Artists
