
'use strict';


module.exports = function (grunt) {
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		version: grunt.file.readJSON('package.json').version,
		banner: '/*! <%= pkg.title || pkg.name %> - v<%= version %>\n' +
			'<%= pkg.homepage ? "* " + pkg.homepage + "\\n" : "" %>' +
			'* Copyright (c) <%= grunt.template.today("yyyy") %> <%= pkg.author.name %>;' +
			' Licensed <%= pkg.license %> */\n',
	});

	grunt.config.set('less', {
		production: {
			files: { 'css/bundle.css': 'less/index.less' },
			options: { compress: true },
			plugins: [
				new (require('less-plugin-autoprefix'))({browsers: ["last 2 versions"]}),
				new (require('less-plugin-clean-css'))({})
			],
		},
	});

	grunt.config.set('watch', {
		options: {
			atBegin: true,
		},
		less: {
			files: ['less/**/*.less'],
			tasks: ['less'],
			options: { spawn: true },
		},
	});

	grunt.loadNpmTasks('grunt-browserify');
	grunt.loadNpmTasks('grunt-contrib-watch');
	grunt.loadNpmTasks('grunt-contrib-less');
	grunt.loadNpmTasks('grunt-contrib-uglify');
	grunt.loadNpmTasks('grunt-nodemon');
	// grunt.loadNpmTasks('grunt-s3');

	grunt.registerTask('serve', ['nodemon:server']);
	grunt.registerTask('build', ['uglify:js']);
	grunt.registerTask('deploy', ['build', 's3:deploy']);
};
