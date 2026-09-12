import { Center, SegmentedControl, Tooltip } from '@mantine/core'
import { IconList, IconListTree } from '@tabler/icons-react'
import { useAppDispatch, useAppSelector } from '../../../../../../../../state/store.ts'
import OutlineView from '../../types/enums/OutlineView.ts'
import { setView } from '../../state/slice/songOutlineSlice.ts'

function SongOutlineViewControl() {
  const dispatch = useAppDispatch()
  const view = useAppSelector((state) => state.songOutline.view)

  function handleChange(view: OutlineView) {
    dispatch(setView(view))
  }

  return (
    <SegmentedControl<OutlineView>
      value={view}
      onChange={handleChange}
      size={'xs'}
      color={'gray'}
      radius={'12px'}
      data={[
        {
          value: OutlineView.Sections,
          label: (
            <Tooltip label={'Sections View'}>
              <Center>
                <IconListTree size={16} aria-label={'sections-view'} />
              </Center>
            </Tooltip>
          )
        },
        {
          value: OutlineView.Parts,
          label: (
            <Tooltip label={'Parts View'}>
              <Center>
                <IconList size={16} aria-label={'parts-view'} />
              </Center>
            </Tooltip>
          )
        }
      ]}
      styles={{
        label: {
          width: '24px',
          padding: '3px'
        }
      }}
    />
  )
}

export default SongOutlineViewControl
