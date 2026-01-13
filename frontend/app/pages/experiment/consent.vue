<template>
    <div class="container mx-auto px-4 md:px-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 lg:gap-16 items-start max-w-6xl mx-auto">

            <section class="space-y-6 lg:sticky lg:top-24">
                <div class="space-y-2">
                    <h1 class="text-3xl md:text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-[hsl(var(--maz-primary))] to-[hsl(var(--maz-info))]">
                        Before We Begin
                    </h1>
                    <p class="text-lg text-[hsl(var(--maz-muted))]">
                        Please review the study protocols and provide your anonymous demographic details.
                    </p>
                </div>

                <div class="bg-[hsl(var(--maz-secondary))] rounded-xl p-6 border border-[hsl(var(--maz-border))] space-y-4 shadow-sm">
                    <h3 class="text-sm font-bold uppercase tracking-wider text-[hsl(var(--maz-primary))] flex items-center gap-2">
                        <MazIcon name="shield-check" class="w-4 h-4" /> 
                        Study Information
                    </h3>
                    <p class="text-sm leading-relaxed opacity-80">
                        You are invited to participate in a research study examining ethical decision-making in autonomous vehicle scenarios. 
                        This study focuses on <strong>priority ranking</strong> rather than binary choices.
                    </p>
                    <p class="text-sm leading-relaxed opacity-80">
                        The simulation consists of 5 interactive scenarios and takes approximately <strong>5-10 minutes</strong>.
                    </p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="p-4 rounded-xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))]/50">
                        <h4 class="font-semibold mb-2 flex items-center gap-2">
                            <MazIcon name="lock-closed" class="w-4 h-4 text-[hsl(var(--maz-info))]" />
                            Your Rights
                        </h4>
                        <ul class="text-xs space-y-2 opacity-70 list-disc list-inside">
                            <li>Participation is voluntary</li>
                            <li>Withdraw at any time</li>
                            <li>Data is strictly confidential</li>
                        </ul>
                    </div>
                    
                    <div class="p-4 rounded-xl border border-[hsl(var(--maz-border))] bg-[hsl(var(--maz-background))]/50">
                        <h4 class="font-semibold mb-2 flex items-center gap-2">
                            <MazIcon name="cpu-chip" class="w-4 h-4 text-[hsl(var(--maz-warning))]" />
                            Data Usage
                        </h4>
                        <ul class="text-xs space-y-2 opacity-70 list-disc list-inside">
                            <li>GDPR Compliant</li>
                            <li>No PII (IP/Email) stored</li>
                            <li>Aggregate analysis only</li>
                        </ul>
                    </div>
                </div>
            </section>

            <section class="bg-[hsl(var(--maz-secondary))] rounded-2xl p-6 md:p-8 border border-[hsl(var(--maz-border))] shadow-xl shadow-[hsl(var(--maz-primary))]/5">
                <form @submit.prevent="handleStart" class="space-y-6">
                    
                    <div class="space-y-4">
                        <h2 class="text-xl font-semibold border-b border-[hsl(var(--maz-border))] pb-2">
                            Demographics
                        </h2>

                        <div class="grid grid-cols-2 gap-4">
                            <MazSelect
                                v-model="form.age_range"
                                label="Age Range"
                                :options="ageRangeOptions"
                                required
                            />
                            <MazSelect
                                v-model="form.gender"
                                label="Gender"
                                :options="genderOptions"
                                required
                            />
                        </div>

                        <MazSelect
                            v-model="form.country"
                            label="Country of Residence"
                            :options="countryOptions"
                            search
                            required
                            placeholder="Select your region"
                        />

                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <MazSelect
                                v-model="form.occupation"
                                label="Field of Work/Study"
                                :options="occupationOptions"
                                required
                            />
                            <MazSelect
                                v-model="form.driving_experience"
                                label="Do you drive?"
                                :options="drivingOptions"
                                required
                            />
                        </div>

                        <div class="pt-2">
                            <label class="text-sm font-medium block mb-2 opacity-90">
                                Have you participated in this experiment before?
                            </label>
                            <div class="flex gap-6">
                                <MazRadio
                                    v-model="form.has_participated"
                                    :value="false"
                                    name="participated"
                                    color="success"
                                >
                                    No, this is my first time
                                </MazRadio>
                                <MazRadio
                                    v-model="form.has_participated"
                                    :value="true"
                                    name="participated"
                                    color="warning"
                                >
                                    Yes, I have done this before
                                </MazRadio>
                            </div>
                            <p class="text-xs text-[hsl(var(--maz-muted))] mt-1 italic">
                                *Being honest helps us categorize the data correctly. You can still participate!
                            </p>
                        </div>
                    </div>

                    <div class="pt-6 border-t border-[hsl(var(--maz-border))]">
                        <MazCheckbox
                            v-model="form.consent"
                            color="primary"
                            size="md"
                        >
                            <span class="text-sm opacity-80">
                                I have read and understood the study information above. I understand that my participation is voluntary.
                            </span>
                        </MazCheckbox>
                    </div>

                    <MazBtn
                        type="submit"
                        size="xl"
                        block
                        :loading="store.isLoading"
                        :disabled="!isValid"
                        class="bg-sensor-gradient hover:shadow-lg hover:shadow-[hsl(var(--maz-primary))]/30 transition-all duration-300"
                    >
                        Start Simulation
                    </MazBtn>

                </form>
            </section>
        </div>
    </div>
</template>

<script setup lang="ts">
// --- Imports ---
import { useExperimentStore } from '~/stores/experiment'
import countries from "i18n-iso-countries"
import enLocale from "i18n-iso-countries/langs/en.json"

// --- Page Meta ---
definePageMeta({
    layout: 'optional-layout',
    middleware: ['session']
})

// --- Composables ---
const store = useExperimentStore()
const router = useRouter()

// --- i18n Setup ---
countries.registerLocale(enLocale)

// --- Static Options ---
const ageRangeOptions = [
    { label: '18-24', value: 1 },
    { label: '25-34', value: 2 },
    { label: '35-44', value: 3 },
    { label: '45-54', value: 4 },
    { label: '55-64', value: 5 },
    { label: '65+', value: 6 }
]

const genderOptions = [
    { label: 'Male', value: 1 },
    { label: 'Female', value: 2 },
    { label: 'Non-binary / Third gender', value: 3 },
    { label: 'Prefer not to say', value: 4 }
]

const occupationOptions = [
    { label: 'Computer Science / AI / Tech', value: 'tech' },
    { label: 'Ethics / Philosophy / Law', value: 'ethics_law' },
    { label: 'Transport / Automotive', value: 'transport' },
    { label: 'Student (Other)', value: 'student_other' },
    { label: 'General Public (Other)', value: 'general' }
]

const drivingOptions = [
    { label: 'Yes, I hold a license', value: 1 },
    { label: 'No, I do not drive', value: 2 },
    { label: 'I am learning', value: 3 }
]

// --- Reactive State ---
const form = reactive({
    age_range: undefined as number | undefined,
    gender: undefined as number | undefined,
    country: undefined as string | undefined,
    occupation: undefined as string | undefined,
    driving_experience: undefined as number | undefined,
    has_participated: false,
    consent: false
})

// --- Computed ---
const countryOptions = computed(() => {
    const countryObj = countries.getNames("en", { select: "official" })
    return Object.entries(countryObj)
        .map(([code, name]) => ({ label: name, value: code }))
        .sort((a, b) => a.label.localeCompare(b.label))
})

const isValid = computed(() => {
    return (
        form.age_range &&
        form.gender &&
        form.country &&
        form.occupation &&
        form.driving_experience &&
        form.consent
    )
})

// --- Methods ---
async function handleStart() {
    if (!isValid.value) return

    await store.createSession(
        {
            age_range: form.age_range,
            gender: form.gender,
            country: form.country,
            occupation: form.occupation,
            driving_experience: form.driving_experience
        },
        !form.has_participated
    )

    if (store.token) {
        // First time users go to guide, returning users go directly to experiment
        if (!form.has_participated) {
            store.setGuideAccess(true)
            router.push('/guide')
        } else {
            store.setGuideAccess(false)
            router.push('/experiment')
        }
    }
}
</script>