import { Center, Group, Stack, Textarea, TextInput, Tooltip } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import ArtistAutocomplete from '../../@ui/form/input/ArtistAutocomplete.tsx'
import AlbumAutocomplete from '../../@ui/form/input/AlbumAutocomplete.tsx'
import { AddNewSongForm } from '../../../validation/songsForm.ts'
import { IconInfoCircleFilled } from '@tabler/icons-react'
import { useDidUpdate } from '@mantine/hooks'
import { AlbumSearch, ArtistSearch } from '../../../types/models/Search.ts'

interface AddNewSongModalFirstStepProps {
  form: UseFormReturnType<AddNewSongForm>
  artist: ArtistSearch
  setArtist: (artist: ArtistSearch) => void
  album: AlbumSearch
  setAlbum: (album: AlbumSearch) => void
}

function AddNewSongModalFirstStep({
  form,
  artist,
  setArtist,
  album,
  setAlbum
}: AddNewSongModalFirstStepProps) {
  useDidUpdate(() => {
    setArtist(album?.artist as ArtistSearch)
    form.setFieldValue('artistName', album?.artist?.name)
  }, [album])

  return (
    <Stack>
      <TextInput
        withAsterisk={true}
        maxLength={100}
        label="Title"
        placeholder="The title of the song"
        key={form.key('title')}
        {...form.getInputProps('title')}
      />

      <Group>
        <AlbumAutocomplete
          album={album}
          setAlbum={setAlbum}
          key={form.key('albumTitle')}
          setValue={(v) => form.setFieldValue('albumTitle', v)}
          {...form.getInputProps('albumTitle')}
        />

        <Group flex={1} gap={0}>
          <ArtistAutocomplete
            artist={artist}
            setArtist={setArtist}
            key={form.key('artistName')}
            setValue={(v) => form.setFieldValue('artistName', v)}
            {...form.getInputProps('artistName')}
          />
          {album && (
            <Center c={'primary.8'} mt={'lg'} ml={4}>
              <Tooltip
                multiline
                w={210}
                ta={'center'}
                label={'Song will inherit artist from album (even if it has one or not)'}
              >
                <IconInfoCircleFilled aria-label={'artist-info'} size={18} />
              </Tooltip>
            </Center>
          )}
        </Group>
      </Group>

      <Textarea
        label="Description"
        placeholder="Enter Description"
        autosize={true}
        minRows={3}
        maxRows={6}
        key={form.key('description')}
        {...form.getInputProps('description')}
      />
    </Stack>
  )
}

export default AddNewSongModalFirstStep
