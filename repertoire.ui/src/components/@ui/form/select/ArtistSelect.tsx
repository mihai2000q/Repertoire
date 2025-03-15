import {
  Avatar,
  Combobox,
  Group,
  HoverCard,
  LoadingOverlay,
  ScrollArea,
  Text,
  TextInput,
  TextInputProps,
  useCombobox
} from '@mantine/core'
import { IconUserFilled } from '@tabler/icons-react'
import { useEffect, useState } from 'react'
import artistPlaceholder from '../../../../assets/user-placeholder.jpg'
import { useDebouncedValue } from '@mantine/hooks'
import { useGetSearchQuery } from '../../../../state/api/searchApi.ts'
import SearchType from '../../../../utils/enums/SearchType.ts'
import { ArtistSearch } from '../../../../types/models/Search.ts'

interface ArtistSelectProps extends TextInputProps {
  artist: ArtistSearch | null
  setArtist: (artist: ArtistSearch | null) => void
}

function ArtistSelect({ artist, setArtist, ...others }: ArtistSelectProps) {
  const combobox = useCombobox({
    onDropdownClose: () => combobox.resetSelectedOption()
  })

  const [value, setValue] = useState('')
  const [search, setSearch] = useState('')
  const [searchQuery] = useDebouncedValue(search, 200)

  useEffect(() => {
    setValue(artist?.name ?? '')
    setSearch(artist?.name ?? '')
  }, [artist])

  const { data: artists, isFetching } = useGetSearchQuery({
    query: searchQuery,
    currentPage: 1,
    pageSize: 10,
    type: SearchType.Artist,
    order: ['updatedAt:desc']
  })

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

  const ArtistOption = ({ localArtist }: { localArtist: ArtistSearch }) => (
    <Combobox.Option
      key={localArtist.id}
      value={localArtist.name}
      aria-label={localArtist.name}
      onClick={() => setArtist(artist?.id === localArtist?.id ? null : localArtist)}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <Avatar
          size={'sm'}
          src={localArtist.imageUrl ?? artistPlaceholder}
          alt={localArtist.name}
        />
        <Text inline fw={500} lineClamp={2}>
          {localArtist.name}
        </Text>
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
    setArtist(null)
  }

  return (
    <Combobox onOptionSubmit={handleSubmit} store={combobox}>
      <Combobox.Target>
        <TextInput
          flex={1}
          maxLength={100}
          label={'Artist'}
          placeholder={'Choose an artist'}
          leftSection={artist ? <ArtistHoverCard /> : <IconUserFilled size={20} />}
          rightSection={
            artist && others.disabled !== true ? (
              <Combobox.ClearButton onClear={handleClear} />
            ) : (
              <Combobox.Chevron />
            )
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
            {artists?.totalCount === 0 && search.trim() === '' ? (
              <Combobox.Empty>There are no artists</Combobox.Empty>
            ) : artists?.totalCount === 0 ? (
              <Combobox.Empty>No artists found</Combobox.Empty>
            ) : (
              artists?.models.map((artist) => (
                <ArtistOption key={artist.id} localArtist={artist as ArtistSearch} />
              ))
            )}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default ArtistSelect
