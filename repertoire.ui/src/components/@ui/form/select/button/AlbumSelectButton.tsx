import {
  ActionIcon,
  ActionIconProps,
  Avatar,
  Center,
  Combobox,
  Group,
  LoadingOverlay,
  ScrollArea,
  Stack,
  Text,
  Tooltip,
  useCombobox
} from '@mantine/core'
import { IconDisc, IconSearch } from '@tabler/icons-react'
import { forwardRef, ReactNode, useEffect, useState } from 'react'
import { useDebouncedValue, useInputState } from '@mantine/hooks'
import { useGetSearchQuery } from '../../../../../state/api/searchApi.ts'
import SearchType from '../../../../../types/enums/SearchType.ts'
import { AlbumSearch } from '../../../../../types/models/Search.ts'
import CustomIconAlbumVinyl from '../../../icons/CustomIconAlbumVinyl.tsx'

interface AlbumSelectButtonProps extends ActionIconProps {
  album: AlbumSearch | null
  setAlbum: (album: AlbumSearch | null) => void
  searchFilter?: string[]
  icon?: ReactNode
}

const AlbumSelectButton = forwardRef<HTMLButtonElement, AlbumSelectButtonProps>(
  ({ album, setAlbum, searchFilter, icon, ...others }, ref) => {
    const [value, setValue] = useState<string>(album?.title ?? '')
    const [search, setSearch] = useInputState(album?.title ?? '')
    const [searchQuery] = useDebouncedValue(search, 200)

    const combobox = useCombobox({
      onDropdownClose: () => {
        combobox.resetSelectedOption()
        combobox.focusTarget()
        setSearch(value)
      },
      onDropdownOpen: () => {
        combobox.focusSearchInput()
      }
    })

    const {
      data: albums,
      isFetching,
      isLoading
    } = useGetSearchQuery({
      query: searchQuery,
      currentPage: 1,
      pageSize: 10,
      type: SearchType.Album,
      order: ['updatedAt:desc'],
      filter: searchFilter
    })

    useEffect(() => {
      setValue(album?.title ?? '')
      setSearch(album?.title ?? '')
    }, [album])

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
            alt={localAlbum.imageUrl && localAlbum.title}
            bg={'gray.5'}
          >
            <Center c={'white'}>
              <CustomIconAlbumVinyl size={12} />
            </Center>
          </Avatar>
          <Stack gap={0}>
            <Text lh={'xxs'} fw={500} lineClamp={localAlbum.artist ? 1 : 2}>
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
      setAlbum(null)
    }

    return (
      <Combobox onOptionSubmit={handleSubmit} store={combobox} withArrow>
        <Combobox.Target withAriaAttributes={false}>
          {album ? (
            <Tooltip
              label={
                <Text fz={'sm'} c={'white'} lineClamp={2}>
                  {album.title} is selected
                </Text>
              }
              openDelay={200}
              multiline={true}
              maw={250}
              disabled={combobox.dropdownOpened}
            >
              <ActionIcon
                ref={ref}
                variant={'transparent'}
                size={'md'}
                radius={'50%'}
                aria-label={album.title}
                aria-selected={true}
                onClick={() => combobox.toggleDropdown()}
                {...others}
              >
                <Avatar
                  src={album.imageUrl}
                  alt={album.imageUrl && album.title}
                  bg={'gray.5'}
                  sx={{ transition: 'filter 0.25s', '&:hover': { filter: 'brightness(0.7)' } }}
                >
                  <Center c={'white'}>
                    <CustomIconAlbumVinyl size={10} />
                  </Center>
                </Avatar>
              </ActionIcon>
            </Tooltip>
          ) : (
            <Tooltip disabled={combobox.dropdownOpened} label={'Select an album'} openDelay={500}>
              <ActionIcon
                ref={ref}
                variant={'form'}
                aria-label={'album'}
                aria-selected={false}
                disabled={isLoading}
                onClick={() => combobox.toggleDropdown()}
                {...others}
              >
                {icon ?? <IconDisc size={18} />}
              </ActionIcon>
            </Tooltip>
          )}
        </Combobox.Target>

        <Combobox.Dropdown miw={180} pt={2} px={'xxs'} pb={0}>
          <Combobox.Search
            size={'xs'}
            px={'xxs'}
            pos={'relative'}
            maxLength={100}
            aria-label={'search'}
            placeholder={'Search'}
            leftSection={<IconSearch size={12} />}
            rightSection={album && <Combobox.ClearButton onClear={handleClear} />}
            value={search}
            onChange={setSearch}
            sx={{
              '.mantine-Input-section': {
                position: 'absolute',
                pointerEvents: 'all'
              },
              '.mantine-Combobox-input': { paddingTop: 1 }
            }}
          />
          <Combobox.Options>
            <LoadingOverlay visible={isFetching} />

            <ScrollArea.Autosize mah={200} scrollbarSize={5}>
              <Stack gap={0} pb={'xxs'}>
                {albums?.totalCount === 0 && search.trim() === '' ? (
                  <Combobox.Empty>There are no albums</Combobox.Empty>
                ) : albums?.totalCount === 0 ? (
                  <Combobox.Empty>No albums found</Combobox.Empty>
                ) : (
                  albums?.models.map((album) => (
                    <AlbumOption key={album.id} localAlbum={album as AlbumSearch} />
                  ))
                )}
              </Stack>
            </ScrollArea.Autosize>
          </Combobox.Options>
        </Combobox.Dropdown>
      </Combobox>
    )
  }
)

AlbumSelectButton.displayName = 'AlbumSelectButton'

export default AlbumSelectButton
