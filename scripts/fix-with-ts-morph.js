// Advanced TypeScript errors fix script using ts-morph
const { Project } = require('ts-morph');
const path = require('path');
const fs = require('fs');

console.log('Starting ts-morph TypeScript fixes...');

// Initialize project
const project = new Project({
  tsConfigFilePath: path.join(__dirname, '../frontend/tsconfig.json'),
  skipAddingFilesFromTsConfig: true
});

// Add source files
console.log('Adding Vue and TS files to project...');
project.addSourceFilesAtPaths([
  '../frontend/src/**/*.ts',
  '../frontend/src/**/*.vue'
]);

// Function to fix snake_case to camelCase in all files
function fixNamingConventions() {
  console.log('Fixing naming conventions (snake_case to camelCase)...');
  
  const sourceFiles = project.getSourceFiles();
  
  const snakeToCamelMap = {
    'image_url': 'imageUrl',
    'product_id': 'productId',
    'product_type': 'productType',
    'item_count': 'itemCount',
    'updated_at': 'updatedAt',
    'created_at': 'createdAt',
    'tracking_number': 'trackingNumber',
    'custom_fields': 'customFields',
    'is_public': 'isPublic',
    'device_name': 'deviceName',
    'device_brand': 'deviceBrand',
    'is_favorite': 'isFavorite',
    'is_wishlist': 'isWishlist'
  };
  
  sourceFiles.forEach(sourceFile => {
    let fileText = sourceFile.getFullText();
    let modified = false;
    
    // Replace all snake_case with camelCase
    for (const [snakeCase, camelCase] of Object.entries(snakeToCamelMap)) {
      if (fileText.includes(snakeCase)) {
        fileText = fileText.replace(new RegExp(snakeCase, 'g'), camelCase);
        modified = true;
      }
    }
    
    if (modified) {
      sourceFile.replaceWithText(fileText);
      console.log(`Fixed naming in: ${sourceFile.getFilePath()}`);
    }
  });
}

// Function to add missing type imports
function fixMissingImports() {
  console.log('Fixing missing type imports...');
  
  // Find files with DeviceListResponse usage but no import
  const deviceListResponseFiles = project.getSourceFiles().filter(
    file => file.getFullText().includes('DeviceListResponse') && 
           !file.getImportDeclarations().some(importDecl => 
             importDecl.getModuleSpecifierValue().includes('device') && 
             importDecl.getNamedImports().some(named => named.getName() === 'DeviceListResponse')
           )
  );
  
  deviceListResponseFiles.forEach(file => {
    if (file.getFilePath().endsWith('.vue') || file.getFilePath().endsWith('.ts')) {
      let fileText = file.getFullText();
      
      // For Vue files, add import after <script> tag
      if (file.getFilePath().endsWith('.vue')) {
        fileText = fileText.replace(
          /<script.*?>/,
          match => `${match}\nimport { DeviceListResponse } from "@/api/device";`
        );
      } else {
        // For TS files, add at the top
        fileText = `import { DeviceListResponse } from "@/api/device";\n${fileText}`;
      }
      
      file.replaceWithText(fileText);
      console.log(`Added DeviceListResponse import to: ${file.getFilePath()}`);
    }
  });
  
  // Find files with ElMessage usage but no import
  const elMessageFiles = project.getSourceFiles().filter(
    file => {
      // Skip files that use global window.ElMessage
      if (file.getFullText().includes('window.ElMessage')) {
        return false;
      }
      
      // Check for ElMessage usage
      if (!file.getFullText().includes('ElMessage')) {
        return false;
      }
      
      // Check if it's already imported
      try {
        return !file.getImportDeclarations().some(importDecl => {
          try {
            return importDecl.getModuleSpecifierValue().includes('element-plus') && 
                   importDecl.getNamedImports().some(named => named.getName() === 'ElMessage');
          } catch (error) {
            console.log(`Skipping import in file ${file.getFilePath()} due to error: ${error.message}`);
            return false;
          }
        });
      } catch (error) {
        console.log(`Error processing file ${file.getFilePath()}: ${error.message}`);
        return false;
      }
    }
  );
  
  elMessageFiles.forEach(file => {
    if (file.getFilePath().endsWith('.vue') || file.getFilePath().endsWith('.ts')) {
      let fileText = file.getFullText();
      
      // For Vue files, add import after <script> tag
      if (file.getFilePath().endsWith('.vue')) {
        fileText = fileText.replace(
          /<script.*?>/,
          match => `${match}\nimport { ElMessage } from "element-plus";`
        );
      } else {
        // For TS files, add at the top
        fileText = `import { ElMessage } from "element-plus";\n${fileText}`;
      }
      
      file.replaceWithText(fileText);
      console.log(`Added ElMessage import to: ${file.getFilePath()}`);
    }
  });
}

// Function to fix SVG type casting issues
function fixSvgTypeCasting() {
  console.log('Fixing SVG type casting issues...');
  
  const svgServiceFile = project.getSourceFile(sourceFile => 
    sourceFile.getFilePath().includes('svgService.ts')
  );
  
  if (svgServiceFile) {
    let fileText = svgServiceFile.getFullText();
    
    // Fix SVG element casting
    fileText = fileText.replace(
      /document\.getElementById\(.*?\) as SVGSVGElement/g,
      match => `${match.replace(' as SVGSVGElement', '')} as unknown as SVGSVGElement`
    );
    
    svgServiceFile.replaceWithText(fileText);
    console.log('Fixed SVG type casting in svgService.ts');
  }
}

// Function to fix battery property optionality
function fixBatteryProperty() {
  console.log('Fixing battery property optionality...');
  
  const deviceFormFile = project.getSourceFile(sourceFile => 
    sourceFile.getFilePath().includes('DeviceForm.vue')
  );
  
  if (deviceFormFile) {
    let fileText = deviceFormFile.getFullText();
    
    // Make battery property optional
    fileText = fileText.replace(
      /battery: {/g,
      'battery?: {'
    );
    
    deviceFormFile.replaceWithText(fileText);
    console.log('Fixed battery property in DeviceForm.vue');
  }
}

// Function to fix property value access on numbers
function fixValuePropertyOnNumbers() {
  console.log('Fixing value property access on numbers...');
  
  const mouseShapeVisualizationFile = project.getSourceFile(sourceFile => 
    sourceFile.getFilePath().includes('MouseShapeVisualization.vue')
  );
  
  if (mouseShapeVisualizationFile) {
    let fileText = mouseShapeVisualizationFile.getFullText();
    
    // Remove .value from number types
    fileText = fileText.replace(/(\d+)\.value/g, '$1');
    
    mouseShapeVisualizationFile.replaceWithText(fileText);
    console.log('Fixed value property access in MouseShapeVisualization.vue');
  }
}

// Execute all fixes
fixNamingConventions();
fixMissingImports();
fixSvgTypeCasting();
fixBatteryProperty();
fixValuePropertyOnNumbers();

// Save all changes
console.log('Saving all changes...');
project.saveSync();

console.log('ts-morph fixes completed successfully!');
