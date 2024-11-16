import {
  Avatar,
  Combobox,
  Group,
  HoverCard,
  Loader,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  useCombobox
} from '@mantine/core'
import albumPlaceholder from '../../../assets/image-placeholder-1.jpg'
import { useGetAlbumsQuery } from '../../../state/albumsApi.ts'
import Album from '../../../types/models/Album.ts'
import { ChangeEvent, FocusEvent, useState } from 'react'
import dayjs from 'dayjs'
import { IconDiscFilled } from '@tabler/icons-react'

interface AlbumsComboboxProps {
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

function AlbumAutocomplete({ album, setAlbum, setValue, ...inputProps }: AlbumsComboboxProps) {
  const combobox = useCombobox()

  const [searchValue, setSearchValue] = useState('')

  const { data: albums, isLoading } = useGetAlbumsQuery({})

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
            size={'md'}
            radius={'md'}
            src={album.imageUrl ?? albumPlaceholder}
            alt={album.title}
          />
          <Stack gap={4}>
            <Text inline fw={500} lineClamp={2}>
              {album.title}
            </Text>
            {album.releaseDate && (
              <Text inline fz={'xs'} c={'dimmed'}>
                {dayjs(album.releaseDate).year()}
              </Text>
            )}
          </Stack>
        </Group>
      </HoverCard.Dropdown>
    </HoverCard>
  )

  return isLoading ? (
    <Group gap={'xs'} flex={1}>
      <Loader size={25} />
      <Text fz={'sm'} c={'dimmed'}>
        Loading Albums...
      </Text>
    </Group>
  ) : (
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
          placeholder={`${albums.models.length > 0 ? 'Choose or Create Album' : 'Enter New Album Name'}`}
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

      <Combobox.Dropdown sx={(theme) => ({ boxShadow: theme.shadows.lg })}>
        <Combobox.Options>
          <ScrollArea.Autosize type={'scroll'} mah={200} scrollbarSize={5}>
            {albums.totalCount === 0 ? (
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
                    <Text inline fw={500} lineClamp={2}>
                      {album.title}
                    </Text>
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
