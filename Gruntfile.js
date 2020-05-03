const os = require('os');
const pluginName = 'grafana-csv-plugin';

function getPluginSuffix() {
  if (os.type() === 'Linux') {
    return '_linux_amd64';
  } else if (os.type() === 'Darwin') {
    return '_darwin_amd64';
  } else if (os.type() === 'Windows_NT') {
    return '_windows_amd64.exe';
  } else {
    throw new Error('Unsupported OS: ' + os.type());
  }
}

module.exports = function(grunt) {

  require('load-grunt-tasks')(grunt);

  grunt.loadNpmTasks('grunt-execute');
  grunt.loadNpmTasks('grunt-contrib-clean');

  grunt.initConfig({
    clean: ['dist'],
    shell: {
      gobuild: {
        command: function() {
          return 'go build -o ./dist/' + pluginName + getPluginSuffix() + ' ./pkg';
        },
      },
    },
    copy: {
      src_to_dist: {
        cwd: 'src',
        expand: true,
        src: [
          '**/*',
          '!*.js',
          '!module.js',
          '!**/*.scss',
        ],
        dest: 'dist/',
      },
      pluginDef: {
        expand: true,
        src: ['plugin.json'],
        dest: 'dist/',
      }
    },
    watch: {
      rebuild_all: {
        files: ['src/**/*', 'plugin.json'],
        tasks: ['default'],
        options: {spawn: false},
      },
    },
    babel: {
      options: {
        sourceMap: true,
        presets:  ['es2015'],
        plugins: ['transform-es2015-modules-systemjs', 'transform-es2015-for-of'],
      },
      dist: {
        files: [{
          cwd: 'src',
          expand: true,
          src: [
            '*.js',
            'module.js',
          ],
          dest: 'dist/',
        }]
      },
    },
    sass: {
      options: {
        sourceMap: true,
      },
      dist: {
        files: {},
      },
    },
  });

  grunt.registerTask('default', [
    'clean',
    'copy:src_to_dist',
    'copy:pluginDef',
    'babel',
    'shell:gobuild',
  ]);
};
