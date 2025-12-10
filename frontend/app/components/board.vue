<template>
  <div class="flex flex-col items-center p-10 bg-gray-100">
    
    <div class="inline-block shadow-2xl border-4 border-gray-800">
      <div 
        v-for="row in 11" 
        :key="row" 
        class="flex"
      >
        <div 
          v-for="col in 20" 
          :key="col" 
          class="w-12 h-12 flex items-center justify-center transition-colors duration-100 hover:bg-red-500"
          :class="getCellClass(row, col)"
        >
          </div>
      </div>
    </div>

    <p class="mt-4 font-mono text-gray-600 text-sm">
      Top-down Street Grid Simulation
    </p>
  </div>
</template>

<script setup>
/**
 * LOGIC MAPPING:
 * The image creates layers from top to bottom. We map grid rows to these layers.
 * * Rows 1-2:   Top Buildings (Dark Grey)
 * Row 3:      Top Sidewalk (Light Grey with vertical lines)
 * Rows 4-5:   Top Road Lane (Asphalt)
 * Row 6:      Center Line (Dashed)
 * Rows 7-8:   Bottom Road Lane (Asphalt)
 * Row 9:      Bottom Sidewalk
 * Rows 10-11: Bottom Buildings
 */

const getCellClass = (row, col) => {
  // --- ZONE 1: TOP BUILDINGS (Rows 1-2) ---
  if (row <= 2) {
    // Add a thick border to the bottom of row 2 to simulate the building edge
    const border = row === 2 ? 'border-b-4 border-gray-900' : '';
    return `bg-gray-500 ${border}`;
  }

  // --- ZONE 2: TOP SIDEWALK (Row 3) ---
  if (row === 3) {
    // Use border-r to create the "tiled" look of pavement
    return 'bg-gray-200 border-r border-gray-300 border-b-4 border-gray-300';
  }

  // --- ZONE 3: THE ROAD (Rows 4-8) ---
  
  // The Center Line (Row 6)
  if (row === 6) {
    // Logic: Only paint the cell white if it's every 3rd or 4th column to make it "dashed"
    // We use the road color for the empty spaces
    if (col % 4 === 0) {
      return 'bg-white h-2 my-auto'; // Squish the cell height to look like a line
    }
    return 'bg-gray-800'; // Asphalt
  }

  // Top and Bottom Lanes
  if (row >= 4 && row <= 8) {
    return 'bg-gray-800'; // Pure Asphalt color
  }

  // --- ZONE 4: BOTTOM SIDEWALK (Row 9) ---
  if (row === 9) {
    // Top border to separate from road
    return 'bg-gray-200 border-r border-gray-300 border-t-4 border-gray-300';
  }

  // --- ZONE 5: BOTTOM BUILDINGS (Rows 10-11) ---
  if (row >= 10) {
    // Top border to create depth
    const border = row === 10 ? 'border-t-4 border-gray-900' : '';
    return `bg-gray-500 ${border}`;
  }
  
  return '';
};
</script>