package org.stathissideris.ascii2image.test;

import java.awt.Color;

import com.thoughtworks.xstream.converters.Converter;
import com.thoughtworks.xstream.converters.MarshallingContext;
import com.thoughtworks.xstream.converters.UnmarshallingContext;
import com.thoughtworks.xstream.io.HierarchicalStreamReader;
import com.thoughtworks.xstream.io.HierarchicalStreamWriter;

public class XmlColorConverter implements Converter {

	@Override
	public boolean canConvert(Class clazz) {
		return clazz.getName().equals("java.awt.Color");
	}

	@Override
	public void marshal(Object o, HierarchicalStreamWriter writer, MarshallingContext context) {
		Color c = (Color) o;
		writer.addAttribute("r", Integer.toString(c.getRed()));
		writer.addAttribute("g", Integer.toString(c.getGreen()));
		writer.addAttribute("b", Integer.toString(c.getBlue()));
		writer.addAttribute("a", Integer.toString(c.getAlpha()));
	}

	@Override
	public Object unmarshal(HierarchicalStreamReader reader, UnmarshallingContext context) {
		return new Color(Integer.parseInt(reader.getAttribute("r")),
				Integer.parseInt(reader.getAttribute("g")),
				Integer.parseInt(reader.getAttribute("b")),
				Integer.parseInt(reader.getAttribute("a")));
	}

}
