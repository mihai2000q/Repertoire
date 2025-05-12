import { ActionIcon, Group, Popover, Tooltip } from '@mantine/core'
import PopoverConfirmation from '../../@ui/popover/PopoverConfirmation.tsx'
import BandMemberCompactSelect from '../../@ui/form/select/compact/BandMemberCompactSelect.tsx'
import InstrumentCompactSelect from '../../@ui/form/select/compact/InstrumentCompactSelect.tsx'
import { useState } from 'react'
import { Instrument, SongSection, SongSettings } from '../../../types/models/Song.ts'
import { BandMember } from '../../../types/models/Artist.ts'
import {
  useUpdateAllSongSectionsMutation,
  useUpdateSongSettingsMutation
} from '../../../state/api/songsApi.ts'
import { IconSettings } from '@tabler/icons-react'

interface SongSectionsSettingsPopoverProps {
  settings: SongSettings
  sections: SongSection[]
  songId: string
  bandMembers?: BandMember[]
}

function SongSectionsSettingsPopover({
  sections,
  settings,
  songId,
  bandMembers
}: SongSectionsSettingsPopoverProps) {
  const [updateSettings] = useUpdateSongSettingsMutation()
  const [updateAll, { isLoading: isUpdateAllLoading }] = useUpdateAllSongSectionsMutation()

  const [defaultInstrument, setDefaultInstrument] = useState(settings.defaultInstrument)
  const [defaultBandMember, setDefaultBandMember] = useState(settings.defaultBandMember)

  const [openedSettingsPopover, setOpenedSettingsPopover] = useState(false)
  const [openedUpdatedDefaultInstrumentPopover, setOpenedUpdatedDefaultInstrumentPopover] =
    useState(false)
  const [openedUpdatedDefaultBandMemberPopover, setOpenedUpdatedDefaultBandMemberPopover] =
    useState(false)

  async function handleDefaultInstrumentChange(newInstrument: Instrument | null) {
    setDefaultInstrument(newInstrument)
    await updateSettings({
      settingsId: settings.id,
      defaultInstrumentId: newInstrument?.id,
      defaultBandMemberId: defaultBandMember?.id
    }).unwrap()
    if (newInstrument && sections.filter((s) => s.instrument?.id !== newInstrument.id).length > 0) {
      setOpenedUpdatedDefaultInstrumentPopover(true)
    }
  }

  async function handleDefaultBandMemberChange(newBandMember: BandMember | null) {
    setDefaultBandMember(newBandMember)
    await updateSettings({
      settingsId: settings.id,
      defaultInstrumentId: defaultInstrument?.id,
      defaultBandMemberId: newBandMember?.id
    }).unwrap()
    if (newBandMember && sections.filter((s) => s.bandMember?.id !== newBandMember.id).length > 0) {
      setOpenedUpdatedDefaultBandMemberPopover(true)
    }
  }

  async function handleUpdateAllSectionsInstruments() {
    await updateAll({
      songId: songId,
      instrumentId: defaultInstrument?.id
    }).unwrap()
    setOpenedUpdatedDefaultInstrumentPopover(false)
  }

  async function handleUpdateAllSectionsBandMembers() {
    await updateAll({
      songId: songId,
      bandMemberId: defaultBandMember?.id
    }).unwrap()
    setOpenedUpdatedDefaultBandMemberPopover(false)
  }

  return (
    <Popover
      opened={openedSettingsPopover}
      onChange={setOpenedSettingsPopover}
      transitionProps={{ transition: 'fade-up' }}
      position={'top'}
      closeOnClickOutside={
        !(openedUpdatedDefaultInstrumentPopover || openedUpdatedDefaultBandMemberPopover)
      }
    >
      <Popover.Target>
        <Tooltip label={'Edit settings'} disabled={openedSettingsPopover}>
          <ActionIcon
            aria-label={'settings'}
            variant={'grey'}
            size={'sm'}
            onClick={() => setOpenedSettingsPopover(!openedSettingsPopover)}
          >
            <IconSettings size={16} />
          </ActionIcon>
        </Tooltip>
      </Popover.Target>

      <Popover.Dropdown>
        <Group gap={'xs'}>
          <PopoverConfirmation
            label={"Would you like to update all sections' band members?"}
            popoverProps={{
              opened: openedUpdatedDefaultBandMemberPopover,
              onChange: setOpenedUpdatedDefaultBandMemberPopover,
              transitionProps: { transition: 'skew-down' },
              closeOnClickOutside: !isUpdateAllLoading,
              withinPortal: false
            }}
            isLoading={isUpdateAllLoading}
            onCancel={() => setOpenedUpdatedDefaultBandMemberPopover(false)}
            onConfirm={handleUpdateAllSectionsBandMembers}
          >
            <BandMemberCompactSelect
              bandMember={defaultBandMember}
              setBandMember={handleDefaultBandMemberChange}
              bandMembers={bandMembers}
              withinPortal={false}
              tooltipLabel={'Choose a default band member'}
            />
          </PopoverConfirmation>

          <PopoverConfirmation
            label={"Would you like to update all sections' instruments?"}
            popoverProps={{
              opened: openedUpdatedDefaultInstrumentPopover,
              onChange: setOpenedUpdatedDefaultInstrumentPopover,
              transitionProps: { transition: 'skew-down' },
              closeOnClickOutside: !isUpdateAllLoading,
              withinPortal: false
            }}
            isLoading={isUpdateAllLoading}
            onCancel={() => setOpenedUpdatedDefaultInstrumentPopover(false)}
            onConfirm={handleUpdateAllSectionsInstruments}
          >
            <InstrumentCompactSelect
              instrument={defaultInstrument}
              setInstrument={handleDefaultInstrumentChange}
              withinPortal={false}
              tooltipLabel={'Choose a default instrument'}
            />
          </PopoverConfirmation>
        </Group>
      </Popover.Dropdown>
    </Popover>
  )
}

export default SongSectionsSettingsPopover
