import { useState } from 'react'
import { useGetArtistsQuery } from '../../state/artistsApi.ts'
import {
  ActionIcon,
  Box,
  Group,
  Loader,
  Pagination,
  Space,
  Stack,
  Text,
  Title
} from '@mantine/core'
import ArtistsLoader from '../../components/artists/loader/ArtistsLoader.tsx'
import ArtistCard from '../../components/artists/card/ArtistCard.tsx'
import AddNewArtistModal from '../../components/artists/modal/AddNewArtistModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconArrowsSort, IconFilterFilled, IconPlus } from '@tabler/icons-react'

function Artists() {
  const [currentPage, setCurrentPage] = useState(1)
  const { data: artists, isLoading } = useGetArtistsQuery({
    pageSize: 20,
    currentPage: currentPage
  })

  const [openedAddNewArtistModal, { open: openAddNewArtistModal, close: closeAddNewArtistModal }] =
    useDisclosure(false)

  return (
    <Stack h={'100%'}>
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

      {artists?.totalCount === 0 && <Text mt={'sm'}>There are no artists yet. Try to add one</Text>}
      <Group gap={'xl'}>
        {isLoading && <ArtistsLoader />}
        {artists?.models.map((artist) => <ArtistCard key={artist.id} artist={artist} />)}
      </Group>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'xs'}>
        {!isLoading ? (
          <Pagination
            data-testid={'artists-pagination'}
            value={currentPage}
            onChange={setCurrentPage}
            total={artists?.totalCount > 0 ? artists?.totalCount / artists?.models.length : 0}
          />
        ) : (
          <Loader size={25} />
        )}
      </Box>
    </Stack>
  )
}

export default Artists
