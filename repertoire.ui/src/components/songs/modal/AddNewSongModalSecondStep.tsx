import {
  ActionIcon,
  alpha,
  Button,
  ComboboxItem,
  Group,
  Loader,
  NumberInput,
  Select,
  Stack,
  Text,
  TextInput
} from '@mantine/core'
import { DatePickerInput } from '@mantine/dates'
import { IconGripVertical, IconMinus } from '@tabler/icons-react'
import { useGetGuitarTuningsQuery, useGetSongSectionTypesQuery } from '../../../state/songsApi.ts'
import { Dispatch, SetStateAction } from 'react'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import { v4 as uuid } from 'uuid'
import { UseFormReturnType } from '@mantine/form'
import { AddNewSongModalSongSection } from './AddNewSongModal.tsx'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import { UseListStateHandlers } from '@mantine/hooks'

interface AddNewSongModalSecondStepProps {
  form: UseFormReturnType<unknown, (values: unknown) => unknown>
  sections: AddNewSongModalSongSection[]
  sectionsHandlers: UseListStateHandlers<AddNewSongModalSongSection>
  guitarTuning: ComboboxItem | null
  setGuitarTuning: Dispatch<SetStateAction<ComboboxItem | null>>
  difficulty: ComboboxItem | null
  setDifficulty: Dispatch<SetStateAction<ComboboxItem | null>>
}

function AddNewSongModalSecondStep({
  form,
  sections,
  sectionsHandlers,
  guitarTuning,
  setGuitarTuning,
  difficulty,
  setDifficulty
}: AddNewSongModalSecondStepProps) {
  const { data: guitarTuningsData, isLoading: isGuitarTuningsLoading } = useGetGuitarTuningsQuery()
  const guitarTunings = guitarTuningsData?.map((guitarTuning) => ({
    value: guitarTuning.id,
    label: guitarTuning.name
  }))

  const difficulties = Object.entries(Difficulty).map(([key, value]) => ({
    value: value,
    label: key
  }))

  const { data: songSectionTypesData } = useGetSongSectionTypesQuery()
  const songSectionTypes = songSectionTypesData?.map((type) => ({
    value: type.id,
    label: type.name
  }))

  const handleAddSection = () => sectionsHandlers.append({ id: uuid(), name: '', type: null })

  function onSectionsDragEnd({ destination, source }) {
    sectionsHandlers.reorder({ from: source.index, to: destination?.index || 0 })
  }

  return (
    <Stack>
      <Group justify={'space-between'} align={'center'}>
        {isGuitarTuningsLoading ? (
          <Group gap={'xs'}>
            <Loader size={25} />
            <Text fz={'sm'} c={'dimmed'}>
              Loading Tunings...
            </Text>
          </Group>
        ) : (
          <Select
            flex={1.25}
            label={'Guitar Tuning'}
            placeholder={'Select Guitar Tuning'}
            data={guitarTunings}
            value={guitarTuning ? guitarTuning.value : null}
            onChange={(_, option) => setGuitarTuning(option)}
            maxDropdownHeight={150}
            clearable
          />
        )}

        <Select
          flex={1}
          label={'Difficulty'}
          placeholder={'Select Difficulty'}
          data={difficulties}
          value={difficulty ? difficulty.value : null}
          onChange={(_, option) => setDifficulty(option)}
          clearable
        />
      </Group>

      <Group justify={'space-between'} align={'center'}>
        <NumberInput
          flex={1}
          label="Bpm"
          placeholder="Enter Bpm"
          key={form.key('bpm')}
          {...form.getInputProps('bpm')}
        />

        <DatePickerInput
          flex={1}
          label={'Release Date'}
          placeholder={'Choose the release date'}
          key={form.key('releaseDate')}
          {...form.getInputProps('releaseDate')}
        />
      </Group>

      <Stack gap={4}>
        <Text fw={500} fz={'sm'}>
          Sections
        </Text>
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
                          provided.draggableProps.style.top = provided.draggableProps.style.top - 36
                        }
                      }
                      return (
                        <Group
                          key={section.id}
                          align={'center'}
                          gap={'xs'}
                          py={'xs'}
                          sx={(theme) => ({
                            borderRadius: '16px',
                            border: `1px solid ${alpha(theme.colors.cyan[9], 0.33)}`,
                            borderWidth: snapshot.isDragging ? '1px' : '0px'
                          })}
                          ref={provided.innerRef}
                          {...provided.draggableProps}
                        >
                          <ActionIcon variant={'subtle'} size={'lg'} {...provided.dragHandleProps}>
                            <IconGripVertical size={20} />
                          </ActionIcon>

                          <Select
                            w={95}
                            placeholder={'Type'}
                            data={songSectionTypes}
                            value={section.type ? section.type.value : null}
                            onChange={(_, option) =>
                              sectionsHandlers.setItem(index, { ...section, type: option })
                            }
                            maxDropdownHeight={150}
                          />

                          <TextInput
                            flex={1}
                            maxLength={30}
                            placeholder={'Name of Section'}
                            value={section.name}
                            onChange={(e) =>
                              sectionsHandlers.setItem(index, { ...section, name: e.target.value })
                            }
                          />

                          <ActionIcon
                            variant={'subtle'}
                            size={'lg'}
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
        <Button style={{ alignSelf: 'start' }} variant={'subtle'} onClick={handleAddSection}>
          Add Section
        </Button>
      </Stack>
    </Stack>
  )
}

export default AddNewSongModalSecondStep
