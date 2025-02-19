import {
  Avatar,
  Combobox,
  Group,
  HoverCard,
  LoadingOverlay,
  ScrollArea, Stack,
  Text,
  TextInput,
  TextInputProps,
  useCombobox
} from '@mantine/core'
import Album from '../../../../types/models/Album.ts'
import { IconUserFilled } from '@tabler/icons-react'
import { useState } from 'react'
import { useGetAlbumsQuery } from '../../../../state/api/albumsApi.ts'
import albumPlaceholder from '../../../../assets/image-placeholder-1.jpg'
import {useDebouncedValue} from "@mantine/hooks";
import dayjs from "dayjs";

interface AlbumSelectProps extends TextInputProps {
  album: Album | null
  setAlbum: (album: Album | null) => void
}

function AlbumSelect({ album, setAlbum, ...others }: AlbumSelectProps) {
  const combobox = useCombobox({
    onDropdownClose: () => combobox.resetSelectedOption()
  })

  const [value, setValue] = useState(album?.title ?? '')
  const [search, setSearch] = useState(album?.title ?? '')
  const [searchQuery] = useDebouncedValue(search, 200)

  const { data: albums, isFetching } = useGetAlbumsQuery({
    currentPage: 1,
    pageSize: 10,
    orderBy: ['title asc'],
    searchBy: searchQuery.trim() !== '' ? [`title ~* '${searchQuery.trim()}'`] : []
  })

  const AlbumHoverCard = () => (
    <HoverCard withArrow={true} openDelay={200} position="bottom" shadow={'md'}>
      <HoverCard.Target>
        <Avatar radius={'md'} size={23} src={album.imageUrl ?? albumPlaceholder} alt={album.title} />
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <Group gap={'xs'} maw={200} wrap={'nowrap'}>
          <Avatar
            size={'lg'}
            radius={'md'}
            src={album.imageUrl ?? albumPlaceholder}
            alt={album.title}
          />
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

  const AlbumOption = ({ localAlbum }: { localAlbum: Album }) => (
    <Combobox.Option
      key={localAlbum.id}
      value={localAlbum.title}
      aria-label={localAlbum.title}
      onClick={() => setAlbum(album === localAlbum ? null : localAlbum)}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <Avatar
          size={'sm'}
          radius={'md'}
          src={localAlbum.imageUrl ?? albumPlaceholder}
          alt={localAlbum.title}
        />
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
          leftSection={album ? <AlbumHoverCard /> : <IconUserFilled size={20} />}
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
              albums?.models.map((album) => <AlbumOption key={album.id} localAlbum={album} />)
            )}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default AlbumSelect
