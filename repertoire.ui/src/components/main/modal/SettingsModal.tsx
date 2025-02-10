import { Modal, Tabs } from '@mantine/core'
import { useState } from 'react'
import SettingsModalAccountTab from './SettingsModalAccountTab.tsx'
import SettingsModalCustomizationTab from './SettingsModalCustomizationTab.tsx'
import User from "../../../types/models/User.ts";

enum SettingsTabs {
  Account = 'account',
  Customization = 'customize'
}

interface SettingsModalProps {
  opened: boolean
  onClose: () => void
  user: User
}

function SettingsModal({ opened, onClose, user }: SettingsModalProps) {
  const [activeTab, setActiveTab] = useState<string>(SettingsTabs.Account)

  return (
    <Modal opened={opened} onClose={onClose} title={'Settings'} size={'lg'}>
      <Modal.Body p={0}>
        <Tabs variant={'default'} value={activeTab} onChange={setActiveTab}>
          <Tabs.List>
            <Tabs.Tab value={SettingsTabs.Account}>Account</Tabs.Tab>
            <Tabs.Tab value={SettingsTabs.Customization}>Customization</Tabs.Tab>
          </Tabs.List>

          <Tabs.Panel value={SettingsTabs.Account}>
            <SettingsModalAccountTab user={user} onCloseSettingsModal={onClose} />
          </Tabs.Panel>
          <Tabs.Panel value={SettingsTabs.Customization}>
            <SettingsModalCustomizationTab />
          </Tabs.Panel>
        </Tabs>
      </Modal.Body>
    </Modal>
  )
}

export default SettingsModal
