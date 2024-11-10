import { Autocomplete, Group, Loader, Stack, Text, Textarea, TextInput } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { useGetAlbumsQuery } from '../../../state/albumsApi.ts'
import { useGetArtistsQuery } from '../../../state/artistsApi.ts'

interface AddNewSongModalFirstStepProps {
  form: UseFormReturnType<unknown, (values: unknown) => unknown>
}

function AddNewSongModalFirstStep({ form }: AddNewSongModalFirstStepProps) {
  const { data: albumsData, isLoading: isAlbumsLoading } = useGetAlbumsQuery({})
  const albums = albumsData?.models?.map((album) => ({
    value: album.id,
    label: album.title
  }))

  const { data: artistsData, isLoading: isArtistsLoading } = useGetArtistsQuery({})
  const artists = artistsData?.models?.map((artist) => ({
    value: artist.id,
    label: artist.name
  }))

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

      <Group align={'center'}>
        {isAlbumsLoading ? (
          <Group gap={'xs'} flex={1}>
            <Loader size={25} />
            <Text fz={'sm'} c={'dimmed'}>
              Loading Albums...
            </Text>
          </Group>
        ) : (
          <Autocomplete
            flex={1}
            data={albums}
            label={'Album'}
            placeholder={`${albums.length > 0 ? 'Choose or Create Album' : 'Enter New Album Title'}`}
          />
        )}

        {isArtistsLoading ? (
          <Group gap={'xs'} flex={1}>
            <Loader size={25} />
            <Text fz={'sm'} c={'dimmed'}>
              Loading Artists...
            </Text>
          </Group>
        ) : (
          <Autocomplete
            flex={1}
            data={artists}
            label={'Artist'}
            placeholder={`${albums.length > 0 ? 'Choose or Create Artist' : 'Enter New Artist Name'}`}
          />
        )}
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
