import { Card, Group, Skeleton, Stack } from '@mantine/core'
import OutlineView from '../types/enums/OutlineView.ts'
import SongSectionsLoader from './SongSectionsLoader.tsx'
import SongPartsLoader from './SongPartsLoader.tsx'

function SongOutlineWidgetLoader({ outlineView }: { outlineView: OutlineView }) {
  return (
    <Card variant={'widget'} aria-label={'outline-loader'} p={0}>
      <Stack gap={0}>
        <Group px={'md'} py={'sm'} gap={'xxs'}>
          <Skeleton w={60} h={15} />

          <Group gap={'xxs'}>
            {Array.from({ length: 6 }).map((_, i) => (
              <Skeleton key={i} w={20} h={20} />
            ))}
            <Skeleton radius={'10px'} w={50} h={25} />
          </Group>
        </Group>

        {outlineView === OutlineView.Sections && <SongSectionsLoader />}
        {outlineView === OutlineView.Parts && <SongPartsLoader />}
      </Stack>
    </Card>
  )
}

export default SongOutlineWidgetLoader
