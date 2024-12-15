import {
  Avatar,
  Combobox,
  Group,
  HoverCard,
  LoadingOverlay,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  useCombobox
} from '@mantine/core'
import albumPlaceholder from '../../../assets/image-placeholder-1.jpg'
import { useGetAlbumsQuery } from '../../../state/albumsApi.ts'
import Album from '../../../types/models/Album.ts'
import { ChangeEvent, FocusEvent } from 'react'
import dayjs from 'dayjs'
import { IconDiscFilled } from '@tabler/icons-react'
import { useDebouncedState } from '@mantine/hooks'

interface AlbumsAutocompleteProps {
  album: Album
  setAlbum: (album: Album) => void
  setValue: (value: string) => void
  value?: string
  defaultValue?: string
  error?: string
  onChange?: (event: ChangeEvent<HTMLInputElement>) => void
  onFocus?: (event: FocusEvent<HTMLInputElement>) => void
  onBlur?: (event: FocusEvent<HTMLInputElement>) => void
}

function AlbumAutocomplete({ album, setAlbum, setValue, ...inputProps }: AlbumsAutocompleteProps) {
  const combobox = useCombobox()

  const [searchValue, setSearchValue] = useDebouncedState('', 200)

  const { data: albums, isFetching } = useGetAlbumsQuery({
    currentPage: 1,
    pageSize: 10,
    orderBy: ['title asc'],
    searchBy:
      searchValue.trim() !== ''
        ? [`title ~* '${searchValue}'`]
        : album
          ? [`title ~* '${album.title}'`]
          : []
  })

  function handleClear() {
    if (setValue) setValue('')
    setSearchValue('')
    setAlbum(null)
  }

  const AlbumHoverCard = () => (
    <HoverCard withArrow={true} openDelay={200} position="bottom" shadow={'md'}>
      <HoverCard.Target>
        <Avatar
          radius={'md'}
          size={23}
          src={album.imageUrl ?? albumPlaceholder}
          alt={album.title}
        />
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <Group gap={'xs'} maw={200} align={'center'} wrap={'nowrap'}>
          <Avatar
            size={'lg'}
            radius={'md'}
            src={album.imageUrl ?? albumPlaceholder}
            alt={album.title}
          />
          <Stack gap={4}>
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
                {dayjs(album.releaseDate).format('DD MMM YYYY')}
              </Text>
            )}
          </Stack>
        </Group>
      </HoverCard.Dropdown>
    </HoverCard>
  )

  return (
    <Combobox
      onOptionSubmit={(optionValue) => {
        if (setValue) setValue(optionValue)
        setSearchValue(optionValue)

        combobox.closeDropdown()
      }}
      store={combobox}
    >
      <Combobox.Target>
        <TextInput
          flex={1}
          maxLength={100}
          label={'Album'}
          placeholder={`${albums?.models?.length > 0 ? 'Choose or Create Album' : 'Enter New Album Name'}`}
          leftSection={album ? <AlbumHoverCard /> : <IconDiscFilled size={20} />}
          rightSection={album && <Combobox.ClearButton onClear={handleClear} />}
          {...inputProps}
          onChange={(event) => {
            if (inputProps.onChange) inputProps.onChange(event)
            if (setValue) setValue(event.currentTarget.value)
            setSearchValue(event.currentTarget.value)
            setAlbum(null)

            combobox.openDropdown()
            combobox.updateSelectedOptionIndex()
          }}
          onFocus={(e) => {
            combobox.openDropdown()
            if (inputProps.onFocus) inputProps.onFocus(e)
          }}
          onBlur={(e) => {
            combobox.closeDropdown()
            if (inputProps.onBlur) inputProps.onBlur(e)
          }}
        />
      </Combobox.Target>

      <Combobox.Dropdown style={(theme) => ({ boxShadow: theme.shadows.lg })}>
        <LoadingOverlay visible={isFetching} />

        <Combobox.Options>
          <ScrollArea.Autosize type={'scroll'} mah={200} scrollbarSize={5}>
            {albums?.totalCount === 0 ? (
              <Combobox.Empty>No album found</Combobox.Empty>
            ) : (
              albums?.models?.map((album) => (
                <Combobox.Option value={album.title} key={album.id} onClick={() => setAlbum(album)}>
                  <Group gap={'xs'} align={'center'} wrap={'nowrap'}>
                    <Avatar
                      size={'sm'}
                      radius={'md'}
                      src={album.imageUrl ?? albumPlaceholder}
                      alt={album.title}
                    />
                    <Stack gap={0}>
                      <Text inline fw={500} lineClamp={2}>
                        {album.title}
                      </Text>
                      {album.artist && <Text inline c={'dimmed'} fz={'xs'} fw={500} truncate={'end'}>
                        {album.artist.name}
                      </Text>}
                    </Stack>
                  </Group>
                </Combobox.Option>
              ))
            )}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default AlbumAutocomplete
