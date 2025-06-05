import {
  ActionIcon,
  alpha,
  Button,
  Center,
  ComboboxItem,
  Group,
  NumberInput,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import { DatePickerInput } from '@mantine/dates'
import {
  IconCalendarRepeat,
  IconGripVertical,
  IconInfoCircleFilled,
  IconMinus
} from '@tabler/icons-react'
import { Dispatch, SetStateAction } from 'react'
import { v4 as uuid } from 'uuid'
import { UseFormReturnType } from '@mantine/form'
import { AddNewSongModalSongSection } from './AddNewSongModal.tsx'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import { UseListStateHandlers } from '@mantine/hooks'
import GuitarTuningSelect from '../../@ui/form/select/GuitarTuningSelect.tsx'
import DifficultySelect from '../../@ui/form/select/DifficultySelect.tsx'
import { AddNewSongForm } from '../../../validation/songsForm.ts'
import CustomIconMetronome from '../../@ui/icons/CustomIconMetronome.tsx'
import SongSectionTypeSelect from '../../@ui/form/select/compact/SongSectionTypeSelect.tsx'
import { AlbumSearch } from '../../../types/models/Search.ts'

interface AddNewSongModalSecondStepProps {
  form: UseFormReturnType<AddNewSongForm>
  sections: AddNewSongModalSongSection[]
  sectionsHandlers: UseListStateHandlers<AddNewSongModalSongSection>
  guitarTuning: ComboboxItem | null
  setGuitarTuning: Dispatch<SetStateAction<ComboboxItem | null>>
  difficulty: ComboboxItem | null
  setDifficulty: Dispatch<SetStateAction<ComboboxItem | null>>
  album: AlbumSearch | null
}

function AddNewSongModalSecondStep({
  form,
  sections,
  sectionsHandlers,
  guitarTuning,
  setGuitarTuning,
  difficulty,
  setDifficulty,
  album
}: AddNewSongModalSecondStepProps) {
  const handleAddSection = () =>
    sectionsHandlers.append({ id: uuid(), name: '', type: null, errors: [] })

  function onSectionsDragEnd({ destination, source }) {
    sectionsHandlers.reorder({ from: source.index, to: destination?.index || 0 })
  }

  return (
    <Stack>
      <Group justify={'space-between'}>
        <GuitarTuningSelect option={guitarTuning} onChange={setGuitarTuning} />

        <DifficultySelect option={difficulty} onChange={setDifficulty} />
      </Group>

      <Group justify={'space-between'}>
        <NumberInput
          flex={0.75}
          min={1}
          allowNegative={false}
          allowDecimal={false}
          leftSection={<CustomIconMetronome size={20} />}
          label="Bpm"
          placeholder="Enter Bpm"
          key={form.key('bpm')}
          {...form.getInputProps('bpm')}
        />

        <Group flex={1} gap={'xxs'}>
          <DatePickerInput
            flex={1}
            leftSection={<IconCalendarRepeat size={20} />}
            label={'Release Date'}
            placeholder={'Choose the release date'}
            key={form.key('releaseDate')}
            {...form.getInputProps('releaseDate')}
          />
          {album?.releaseDate && (
            <Center c={'primary.8'} mt={'lg'} ml={4}>
              <Tooltip
                multiline
                w={210}
                ta={'center'}
                label={'If the release date is not set, then it will be inherited from the album'}
              >
                <IconInfoCircleFilled aria-label={'release-date-info'} size={18} />
              </Tooltip>
            </Center>
          )}
        </Group>
      </Group>

      <Stack gap={'xxs'}>
        <Text fw={500} fz={'sm'}>
          Sections
        </Text>
        <ScrollArea.Autosize mah={'35vh'} offsetScrollbars={'y'} scrollbars={'y'} scrollbarSize={7}>
          <DragDropContext onDragEnd={onSectionsDragEnd}>
            <Droppable droppableId="dnd-list" direction="vertical">
              {(provided) => (
                <Stack gap={0} ref={provided.innerRef} {...provided.droppableProps}>
                  {sections.map((section, index) => (
                    <Draggable key={section.id} index={index} draggableId={section.id}>
                      {(provided, snapshot) => {
                        if (snapshot.isDragging) {
                          if ('left' in provided.draggableProps.style) {
                            provided.draggableProps.style.left = 24
                          }
                          if ('top' in provided.draggableProps.style) {
                            provided.draggableProps.style.top -= 36
                          }
                        }
                        return (
                          <Group
                            key={section.id}
                            gap={'xs'}
                            py={'xs'}
                            sx={(theme) => ({
                              transition: '0.25s',
                              borderRadius: '16px',
                              border: `1px solid ${alpha(theme.colors.primary[9], 0.33)}`,
                              borderWidth: snapshot.isDragging ? '1px' : '0px'
                            })}
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                          >
                            <ActionIcon
                              aria-label={'drag-handle'}
                              variant={'subtle'}
                              size={'lg'}
                              {...provided.dragHandleProps}
                            >
                              <IconGripVertical size={20} />
                            </ActionIcon>

                            <SongSectionTypeSelect
                              w={100}
                              comboboxProps={{ position: 'bottom-start', width: 125 }}
                              placeholder={'Type'}
                              option={section.type}
                              onOptionChange={(option) =>
                                sectionsHandlers.setItem(index, {
                                  ...section,
                                  type: option,
                                  errors: [
                                    ...section.errors.filter((e) => e.property !== 'type'),
                                    ...(!option ? [{ property: 'type' }] : [])
                                  ]
                                })
                              }
                              error={section.errors.some((e) => e.property === 'type')}
                            />

                            <TextInput
                              flex={1}
                              maxLength={30}
                              aria-label={'name'}
                              placeholder={'Name of Section'}
                              value={section.name}
                              onChange={(e) =>
                                sectionsHandlers.setItem(index, {
                                  ...section,
                                  name: e.target.value,
                                  errors: [
                                    ...section.errors.filter((e) => e.property !== 'name'),
                                    ...(e.target.value.trim() === '' ? [{ property: 'name' }] : [])
                                  ]
                                })
                              }
                              error={section.errors.some((e) => e.property === 'name')}
                            />

                            <ActionIcon
                              variant={'subtle'}
                              size={'lg'}
                              aria-label={'remove-section'}
                              onClick={() => sectionsHandlers.remove(index)}
                            >
                              <IconMinus size={20} />
                            </ActionIcon>
                          </Group>
                        )
                      }}
                    </Draggable>
                  ))}
                  {provided.placeholder}
                </Stack>
              )}
            </Droppable>
          </DragDropContext>
        </ScrollArea.Autosize>

        <Button style={{ alignSelf: 'start' }} variant={'subtle'} onClick={handleAddSection}>
          Add Section
        </Button>
      </Stack>
    </Stack>
  )
}

export default AddNewSongModalSecondStep
