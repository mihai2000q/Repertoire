import {
  Avatar,
  Center,
  Combobox,
  Group,
  HoverCard,
  LoadingOverlay,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  TextInputProps,
  useCombobox
} from '@mantine/core'
import { IconDiscFilled } from '@tabler/icons-react'
import { useEffect, useState } from 'react'
import { useDebouncedValue } from '@mantine/hooks'
import dayjs from 'dayjs'
import { AlbumSearch } from '../../../../types/models/Search.ts'
import { useGetSearchQuery } from '../../../../state/api/searchApi.ts'
import SearchType from '../../../../utils/enums/SearchType.ts'
import CustomIconAlbumVinyl from '../../icons/CustomIconAlbumVinyl.tsx'

interface AlbumSelectProps extends TextInputProps {
  album: AlbumSearch | null
  setAlbum: (album: AlbumSearch | null) => void
}

function AlbumSelect({ album, setAlbum, ...others }: AlbumSelectProps) {
  const combobox = useCombobox({
    onDropdownClose: () => combobox.resetSelectedOption()
  })

  const [value, setValue] = useState('')
  const [search, setSearch] = useState('')
  const [searchQuery] = useDebouncedValue(search, 200)

  useEffect(() => {
    setSearch(album?.title ?? '')
    setValue(album?.title ?? '')
  }, [album])

  const { data: albums, isFetching } = useGetSearchQuery({
    query: searchQuery,
    currentPage: 1,
    pageSize: 10,
    type: SearchType.Album,
    order: ['updatedAt:desc']
  })

  const AlbumHoverCard = () => (
    <HoverCard withArrow={true} openDelay={200} position="bottom" shadow={'md'}>
      <HoverCard.Target>
        <Avatar radius={'md'} size={23} src={album.imageUrl} alt={album.title} bg={'gray.5'}>
          <Center c={'white'}>
            <CustomIconAlbumVinyl size={11} />
          </Center>
        </Avatar>
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <Group gap={'xs'} maw={200} wrap={'nowrap'}>
          <Avatar radius={'md'} size={'lg'} src={album.imageUrl} alt={album.title} bg={'gray.5'}>
            <Center c={'white'}>
              <CustomIconAlbumVinyl size={25} />
            </Center>
          </Avatar>
          <Stack gap={'xxs'}>
            <Text inline fw={500} lineClamp={2}>
              {album.title}
            </Text>
            {album.artist && (
              <Text inline fw={500} fz={'xs'} c={'dimmed'}>
                {album.artist.name}
              </Text>
            )}
            {album.releaseDate && (
              <Text inline fz={'xxs'} c={'dimmed'}>
                {dayjs(album.releaseDate).format('D MMM YYYY')}
              </Text>
            )}
          </Stack>
        </Group>
      </HoverCard.Dropdown>
    </HoverCard>
  )

  const AlbumOption = ({ localAlbum }: { localAlbum: AlbumSearch }) => (
    <Combobox.Option
      key={localAlbum.id}
      value={localAlbum.title}
      aria-label={localAlbum.title}
      onClick={() => setAlbum(album?.id === localAlbum?.id ? null : localAlbum)}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <Avatar
          radius={'md'}
          size={'sm'}
          src={localAlbum.imageUrl}
          alt={localAlbum.title}
          bg={'gray.5'}
        >
          <Center c={'white'}>
            <CustomIconAlbumVinyl size={12} />
          </Center>
        </Avatar>
        <Stack gap={0}>
          <Text inline fw={500} lineClamp={2}>
            {localAlbum.title}
          </Text>
          {localAlbum.artist && (
            <Text inline c={'dimmed'} fz={'xs'} fw={500} lineClamp={1}>
              {localAlbum.artist.name}
            </Text>
          )}
        </Stack>
      </Group>
    </Combobox.Option>
  )

  function handleSubmit(valueString: string) {
    setValue(valueString)
    setSearch(valueString)
    combobox.closeDropdown()
  }

  function handleClear() {
    setValue('')
    setSearch('')
    setAlbum(null)
  }

  return (
    <Combobox onOptionSubmit={handleSubmit} store={combobox}>
      <Combobox.Target>
        <TextInput
          flex={1}
          maxLength={100}
          label={'Album'}
          placeholder={'Choose an album'}
          leftSection={album ? <AlbumHoverCard /> : <IconDiscFilled size={20} />}
          rightSection={
            album ? <Combobox.ClearButton onClear={handleClear} /> : <Combobox.Chevron />
          }
          value={search}
          onChange={(e) => {
            combobox.openDropdown()
            combobox.updateSelectedOptionIndex()
            setSearch(e.currentTarget.value)
          }}
          onClick={() => combobox.openDropdown()}
          onFocus={() => combobox.openDropdown()}
          onBlur={() => {
            combobox.closeDropdown()
            setSearch(value)
          }}
          {...others}
        />
      </Combobox.Target>

      <Combobox.Dropdown>
        <LoadingOverlay visible={isFetching} />

        <Combobox.Options>
          <ScrollArea.Autosize mah={200} scrollbarSize={5}>
            {albums?.totalCount === 0 && search.trim() === '' ? (
              <Combobox.Empty>There are no albums</Combobox.Empty>
            ) : albums?.totalCount === 0 ? (
              <Combobox.Empty>No albums found</Combobox.Empty>
            ) : (
              albums?.models.map((album) => (
                <AlbumOption key={album.id} localAlbum={album as AlbumSearch} />
              ))
            )}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default AlbumSelect
