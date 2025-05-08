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
  useCombobox
} from '@mantine/core'
import { ChangeEvent, FocusEvent } from 'react'
import dayjs from 'dayjs'
import { IconDiscFilled } from '@tabler/icons-react'
import { useDebouncedState } from '@mantine/hooks'
import { useGetSearchQuery } from '../../../../state/api/searchApi.ts'
import SearchType from '../../../../types/enums/SearchType.ts'
import { AlbumSearch } from '../../../../types/models/Search.ts'
import CustomIconAlbumVinyl from '../../icons/CustomIconAlbumVinyl.tsx'

interface AlbumsAutocompleteProps {
  album: AlbumSearch | null
  setAlbum: (album: AlbumSearch | null) => void
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

  const { data, isFetching } = useGetSearchQuery({
    query: searchValue,
    currentPage: 1,
    pageSize: 10,
    type: SearchType.Album,
    order: ['updatedAt:desc']
  })
  const totalCount = data?.totalCount
  const albums = data?.models as AlbumSearch[]

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
          src={album.imageUrl}
          alt={album.imageUrl && album.title}
          bg={'gray.5'}
        >
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
            <Text lh={'xxs'} fw={500} lineClamp={2}>
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
          placeholder={`${totalCount > 0 ? 'Choose or Create Album' : 'Enter New Album Name'}`}
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

      <Combobox.Dropdown pb={0}>
        <LoadingOverlay visible={isFetching} />

        <Combobox.Options>
          <ScrollArea.Autosize mah={200} scrollbarSize={5}>
            <Stack gap={0} pb={'xxs'}>
              {totalCount === 0 ? (
                <Combobox.Empty>No album found</Combobox.Empty>
              ) : (
                albums?.map((album) => (
                  <Combobox.Option
                    key={album.id}
                    value={album.title}
                    aria-label={album.title}
                    onClick={() => setAlbum(album)}
                  >
                    <Group gap={'xs'} wrap={'nowrap'}>
                      <Avatar
                        radius={'md'}
                        size={'sm'}
                        src={album.imageUrl}
                        alt={album.imageUrl && album.title}
                        bg={'gray.5'}
                      >
                        <Center c={'white'}>
                          <CustomIconAlbumVinyl size={12} />
                        </Center>
                      </Avatar>
                      <Stack gap={0}>
                        <Text lh={'xxs'} fw={500} lineClamp={album.artist ? 1 : 2}>
                          {album.title}
                        </Text>
                        {album.artist && (
                          <Text inline c={'dimmed'} fz={'xs'} fw={500} lineClamp={1}>
                            {album.artist.name}
                          </Text>
                        )}
                      </Stack>
                    </Group>
                  </Combobox.Option>
                ))
              )}
            </Stack>
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default AlbumAutocomplete
