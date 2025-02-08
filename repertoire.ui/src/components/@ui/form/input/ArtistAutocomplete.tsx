import {
  Avatar,
  Combobox,
  Group,
  HoverCard,
  LoadingOverlay,
  ScrollArea,
  Text,
  TextInput,
  useCombobox
} from '@mantine/core'
import artistPlaceholder from '../../../../assets/user-placeholder.jpg'
import { useGetArtistsQuery } from '../../../../state/api/artistsApi.ts'
import Artist from '../../../../types/models/Artist.ts'
import { ChangeEvent, FocusEvent } from 'react'
import { IconUserFilled } from '@tabler/icons-react'
import { useDebouncedState } from '@mantine/hooks'

interface ArtistsAutocompleteProps {
  artist: Artist | null
  setArtist: (artist: Artist | null) => void
  setValue: (value: string) => void
  value?: string
  defaultValue?: string
  error?: string
  onChange?: (event: ChangeEvent<HTMLInputElement>) => void
  onFocus?: (event: FocusEvent<HTMLInputElement>) => void
  onBlur?: (event: FocusEvent<HTMLInputElement>) => void
}

function ArtistAutocomplete({
  artist,
  setArtist,
  setValue,
  ...inputProps
}: ArtistsAutocompleteProps) {
  const combobox = useCombobox()

  const [searchValue, setSearchValue] = useDebouncedState('', 200)

  const { data: artists, isFetching } = useGetArtistsQuery({
    currentPage: 1,
    pageSize: 10,
    orderBy: ['name asc'],
    searchBy:
      searchValue.trim() !== ''
        ? [`name ~* '${searchValue}'`]
        : artist
          ? [`name ~* '${artist.name}'`]
          : []
  })

  function handleClear() {
    if (setValue) setValue('')
    setSearchValue('')
    setArtist(null)
  }

  const ArtistHoverCard = () => (
    <HoverCard withArrow={true} openDelay={200} position="bottom" shadow={'md'}>
      <HoverCard.Target>
        <Avatar size={23} src={artist.imageUrl ?? artistPlaceholder} alt={artist.name} />
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <Group gap={'xs'} maw={200} wrap={'nowrap'}>
          <Avatar size={'md'} src={artist.imageUrl ?? artistPlaceholder} alt={artist.name} />
          <Text inline fw={500} lineClamp={2}>
            {artist.name}
          </Text>
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
          label={'Artist'}
          placeholder={`${artists?.models?.length > 0 ? 'Choose or Create Artist' : 'Enter New Artist Name'}`}
          leftSection={artist ? <ArtistHoverCard /> : <IconUserFilled size={20} />}
          rightSection={artist && <Combobox.ClearButton onClear={handleClear} />}
          onClick={() => combobox.openDropdown()}
          {...inputProps}
          onChange={(event) => {
            if (inputProps.onChange) inputProps.onChange(event)
            if (setValue) setValue(event.currentTarget.value)
            setSearchValue(event.currentTarget.value)
            setArtist(null)

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

      <Combobox.Dropdown>
        <LoadingOverlay visible={isFetching} />

        <Combobox.Options>
          <ScrollArea.Autosize mah={200} scrollbarSize={5}>
            {artists?.totalCount === 0 ? (
              <Combobox.Empty>No artist found</Combobox.Empty>
            ) : (
              artists?.models?.map((artist) => (
                <Combobox.Option
                  key={artist.id}
                  value={artist.name}
                  aria-label={artist.name}
                  onClick={() => setArtist(artist)}
                >
                  <Group gap={'xs'} wrap={'nowrap'}>
                    <Avatar
                      size={'sm'}
                      src={artist.imageUrl ?? artistPlaceholder}
                      alt={artist.name}
                    />
                    <Text inline fw={500} lineClamp={2}>
                      {artist.name}
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

export default ArtistAutocomplete
